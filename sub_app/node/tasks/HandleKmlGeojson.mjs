import 'dotenv/config'
import Queue from "bull"
import fs, { existsSync, mkdirSync } from "fs"
import { exec } from "child_process";

const JOB_NAME = "handle_kml_geojson"

/**
 * 
 * @param {import('nats').NatsConnection} nats 
 * @returns Queue.Queue<any>
 */
function init(nats) {
    const jobQueue = new Queue(JOB_NAME, { redis: { port: process.env.REDIS_PORT, host: process.env.REDIS_HOST, password: process.env.REDIS_PASSWORD, db: 1 } }); // Specify Redis connection using object

    jobQueue.process(function (job, done) {

        try {
            let m_data = job.data;

            if (!existsSync("files")) {
                mkdirSync("files")
            }

            fs.writeFileSync("files/" + m_data.uuid + ".json", JSON.stringify(m_data))

            const childProcess = exec("node scripts/ConvertKMLGeojson.mjs --config=files/" + (m_data.uuid + ".json"))
            // Listen for the output of the command
            childProcess.stdout.on('data', (data) => {
                console.log(`stdout: ${data}`); // Log standard output
                nats.publish(m_data.uuid + "_process", JSON.stringify({
                    status: "success",
                    return: "Success"
                }))
            });

            childProcess.stderr.on('data', (data) => {
                console.error(`stderr: ${data}`); // Log standard error
                // done(new Error(data), {});
                nats.publish(m_data.uuid + "_process", JSON.stringify({
                    status: "error",
                    return: data
                }))
            });

            childProcess.on('error', (error) => {
                console.error(`Error executing command: ${error}`);
                // done(new Error(error), {});
                nats.publish(m_data.uuid + "_process", JSON.stringify({
                    status: "success",
                    return: error
                }))
            });

            childProcess.on('close', (code, signal) => {
                // console.log(`Command execution completed with code: ${code}`);
                console.log("signal", signal);
                if (signal == null || signal == "SIGTERM") {
                    nats.publish(m_data.uuid + "_done", JSON.stringify({
                        status_code: code,
                        status: "finish",
                        return: {
                            download: process.env.APP_HOST + "/files/" + m_data.uuid + ".json"
                        }
                    }))
                    done(null, {
                        code
                    });
                } else {
                    nats.publish(m_data.uuid + "_done", JSON.stringify({
                        status: "error",
                        return: code
                    }))
                    done(null, {
                        code
                    });
                }

            });
        } catch (error) {
            console.log(error);
            nats.publish(m_data.uuid + "_done", JSON.stringify({
                status: "error",
                return: error
            }))
            done(new Error(error.message))
        }
    })

    return jobQueue;
}

export default init;



