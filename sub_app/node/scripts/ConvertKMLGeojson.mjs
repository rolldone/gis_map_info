import parseKml from "parse-kml";
import superagent from "superagent";
import fs, { existsSync } from "fs"
import ts from "@mapbox/togeojson"
import ddomParser from "xmldom"

function main() {
    // Accessing command-line arguments
    const args = process.argv.slice(2); // Remove the 'node' and 'test.js' from the array
    // Check if '--config' parameter is provided
    const configArg = args.find(arg => arg.startsWith('--config='));

    if (configArg) {
        const configFilePath = configArg.split('=')[1];
        // Read and parse the JSON configuration file
        fs.readFile(configFilePath, 'utf8', (err, data) => {
            // Remove it
            if (existsSync(configFilePath) == true) {
                fs.rmSync(configFilePath)
            }
            if (err) {
                console.error(`Error reading the config file: ${err}`);
                return;
            }

            try {
                const config = JSON.parse(data);
                console.log('Config content:', config);
                // URL of the file you want to download
                const fileUrl = config.download;

                // Path where you want to save the downloaded file
                const filePath = './files/' + config.uuid + ".kml";

                // Make a GET request to download the file
                superagent
                    .get(fileUrl)
                    .end((err, res) => {
                        if (err) {
                            console.error('Error downloading the file:', err);
                        } else {
                            // Save the downloaded file
                            fs.writeFile(filePath, res.body, async (err) => {
                                if (err) {
                                    console.error('Error saving the file:', err);
                                } else {
                                    let [data, err1] = await loadKMl(filePath);
                                    if (err1) {
                                        console.error("loadKML :: ", err1);
                                        return;
                                    }
                                    let jsonFilePath = filePath.replace(".kml", ".json");
                                    fs.writeFileSync(jsonFilePath, JSON.stringify(data));
                                }
                            });
                        }
                    });
            } catch (parseError) {
                console.error('Error parsing the config file:', parseError);
            }
        });


        // ... code to read and parse the config file
    } else {
        console.log('No config file path provided.');
    }
}

/**
 * 
 * @param {string} filePath 
 * @returns {Promise<[Array<any>, CustomError]>}
 */
function loadKMl(filePath) {
    return new Promise((resolve) => {
        parseKml.readKml(filePath)
            .then(function (dataString) {
                try {
                    var kml = new ddomParser.DOMParser().parseFromString(dataString);
                    let geojson_data = ts.kml(kml);
                    let collections = [];
                    let props = {
                        index: 0
                    };
                    for (let a = 0; a < geojson_data.features.length; a++) {
                        let featureItem = geojson_data.features[a]
                        collectionGeometryCollectionType(featureItem.geometry, featureItem.properties, collections, props);
                    }
                    for (let a = 0; a < geojson_data.features.length; a++) {
                        let featureItem = geojson_data.features[a]
                        collectionGeometryCollectionType(featureItem.geometry, featureItem.properties, collections, props);
                    }
                    resolve([collections, null]);
                } catch (error) {
                    console.log("geojson_data = ts.kml(kml) :: ", error)
                    resolve([[], error])
                }
            })
            .catch(function (e) {
                resolve([null, e])
            });
    })
}

/**
 * 
 * @param {Object} geometry 
 * @param {Object} properties 
 * @param {Array<any>} collection 
 * @param {Object} props
 * @returns {Array<any>}
 */
function collectionGeometryCollectionType(geometry, properties, collection = [], props) {
    props.index += 1;
    console.log("Process convert :: ", props.index);
    if (geometry.type == "GeometryCollection") {
        for (let a = 0; a < geometry.geometries.length; a++) {
            collectionGeometryCollectionType(geometry.geometries[a], properties, collection, props);
        }
    } else {
        collection.push({
            type: "Feature",
            geometry: convertPolygonZ_to_Polygon(geometry),
            properties,
        })
    }

    return collection;
}

function convertPolygonZ_to_Polygon(geometry) {
    // Convert Polygonz to 2D Polygon by removing z-coordinate values
    for (let a = 0; a < geometry.coordinates.length; a++) {
        const coordinates2D = geometry.coordinates[a].map(coord => [coord[0], coord[1]]);
        // console.log("vmakfavmkfvm --> ", coordinates2D);
        geometry.coordinates[a] = coordinates2D;
    }
    return geometry;
}
main()