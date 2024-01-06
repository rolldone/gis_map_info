// Import required modules
import 'dotenv/config'
import { connect, StringCodec } from "nats";
import HandleKmlGeoJSONQueue from "./tasks/HandleKmlGeojson.mjs"
import Express from "express";
import path, { dirname } from 'path'
import { fileURLToPath } from 'url';


// Create Express application
const app = Express();

// Define the path to the static assets directory
// Get the directory name of the current ES module file
const __filename = fileURLToPath(import.meta.url);
const publicDirectoryPath = path.join(dirname(__filename), 'files');
console.log("publicDirectoryPath", publicDirectoryPath);

// Configure Express to serve static files from the public directory
app.use("/files", Express.static(publicDirectoryPath));

// Start the Express server
const PORT = process.env.APP_PORT || 3000;
app.listen(PORT, "0.0.0.0", () => {
  console.log(`Server is running on port ${PORT}`);
});

app.get('/', (req, res) => {
  // Serve an HTML page that references static assets
  res.send(`
      <!DOCTYPE html>
      <html lang="en">
      <head>
          <meta charset="UTF-8">
          <title>Static Files</title>
      </head>
      <body>
          <h1>Static Files Example</h1>
          <img src="/files/image.jpg" alt="Sample Image">
          <!-- Reference other static files as needed -->
      </body>
      </html>
  `);
});

// create a codec
const sc = StringCodec();

const NATS_HOST = process.env.NATS_HOST + ":" + process.env.NATS_PORT
const NATS_HOST_2 = process.env.NATS_HOST_2 + ":" + process.env.NATS_PORT
const NATS_HOST_3 = process.env.NATS_HOST_3 + ":" + process.env.NATS_PORT

const servers = [
  { servers: [NATS_HOST, NATS_HOST_2, NATS_HOST_3] },
];

function getRandomNumber(min, max) {
  // Generate a random number between min and max (inclusive)
  return Math.floor(Math.random() * (max - min + 1)) + min;
}


await servers.forEach(async (v) => {
  try {
    const nc = await connect(v);
    const handleKmlGeojson = HandleKmlGeoJSONQueue(nc);

    console.log(`connected to ${nc.getServer()}`);

    const sub = nc.subscribe("convert_kml_geojson");
    (async () => {
      for await (const m of sub) {
        let m_data = sc.decode(m.data);
        m_data = JSON.parse(m_data);
        handleKmlGeojson.add({
          ...m_data
        }, {
          delay: getRandomNumber(1, 10) * 1000
        })
      }
      console.log("subscription convert_kml_geojson closed");
    })();

    setTimeout(() => {
      for (let a = 0; a < 10; a++) {
        let uuid = "9dfaca75-4db8-4df5-9542-219e6588404a_" + a;
        // nc.publish("convert_kml_geojson", sc.encode(JSON.stringify({
        //   "uuid": uuid,
        //   // "callback": "http://gis_map_info_app:8080/api/admin/zone_rdtr/validate/callback",
        //   // "process": "http://gis_map_info_app:8080/api/admin/zone_rdtr/validate/process",
        //   "download": "http://100.114.33.35:8080/api/admin/rdtr_file/assets/9dfaca75-4db8-4df5-9542-219e6588404a.kml"
        // })))

        // let listenUUIDProcess = nc.subscribe(uuid + "_process");
        // (async () => {
        //   for await (const m of listenUUIDProcess) {
        //     let m_data = sc.decode(m.data);
        //     // console.log(`${uuid}_process`,m_data)
        //     m_data = JSON.parse(m_data);
        //   }
        //   console.log(`listenUUIDProcess ${uuid}_process closed`);
        // })();

        // let listenUUIDFinish = nc.subscribe(uuid + "_done");
        // (async () => {
        //   for await (const m of listenUUIDFinish) {
        //     let m_data = sc.decode(m.data);
        //     m_data = JSON.parse(m_data);
        //     listenUUIDProcess.unsubscribe()
        //     listenUUIDFinish.unsubscribe()
        //     // process.exit(0);
        //   }
        //   console.log(`listenUUIDFinish ${uuid}_done closed`);
        // })();

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
    console.log(`error connecting to `, err);
  }
});