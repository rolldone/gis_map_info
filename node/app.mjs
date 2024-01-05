// Import required modules
import { connect, StringCodec } from "nats";
import fs from "fs"
import { exec, execSync } from "child_process";

// create a codec
const sc = StringCodec();

const servers = [
  { servers: ["gis_map_info_nats:4222", "gis_map_info_nats_2:4222", "gis_map_info_nats_3:4222"] },
];

function getRandomNumber(min, max) {
  // Generate a random number between min and max (inclusive)
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

await servers.forEach(async (v) => {
  try {
    const nc = await connect(v);
    console.log(`connected to ${nc.getServer()}`);

    const sub = nc.subscribe("convert_kml_geojson");
    (async () => {
      for await (const m of sub) {
        let m_data = sc.decode(m.data);
        m_data = JSON.parse(m_data);
        fs.writeFileSync(m_data.uuid + ".json", JSON.stringify(m_data))
        const childProcess = exec("node convert_kml_geojson.mjs --config=" + (m_data.uuid + ".json"))
        // Listen for the output of the command
        childProcess.stdout.on('data', (data) => {
          console.log(`stdout: ${data}`); // Log standard output
        });

        childProcess.stderr.on('data', (data) => {
          console.error(`stderr: ${data}`); // Log standard error
        });

        childProcess.on('error', (error) => {
          console.error(`Error executing command: ${error}`);
        });

        childProcess.on('close', (code) => {
          console.log(`Command execution completed with code: ${code}`);
        });
      }
      console.log("subscription convert_kml_geojson closed");
    })();

    setTimeout(() => {
      for (let a = 0; a < 1; a++) {
        nc.publish("convert_kml_geojson", sc.encode(JSON.stringify({
          "uuid": "9dfaca75-4db8-4df5-9542-219e6588404a",
          "callback": "http://gis_map_info_app:8080/api/admin/rdtr_file/validate/callback",
          "data": {
            filColor: "#ffff00"
          },
          "download": "http://100.114.33.35:8080/api/admin/rdtr_files/assets/9dfaca75-4db8-4df5-9542-219e6588404a.kml"
        })))
      }
    }, 1000)

    // this promise indicates the client closed
    // const done = nc.closed();
    // // do something with the connection

    // // close the connection
    // await nc.close();
    // // check if the close was OK
    // const err = await done;
    // if (err) {
    //   console.log(`error closing:`, err);
    // }
  } catch (err) {
    console.log(`error connecting to ${v}`);
  }
});