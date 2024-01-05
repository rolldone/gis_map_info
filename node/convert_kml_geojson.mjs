import parseKml from "parse-kml";
import fs from "fs"

function main() {
    // Accessing command-line arguments
    const args = process.argv.slice(2); // Remove the 'node' and 'test.js' from the array
    console.log("args ,", args)
    // Check if '--config' parameter is provided
    const configIndex = args.indexOf('--config');
    console.log("vvvvvvv", (configIndex == -1))
    console.log("configIndex ", configIndex);
    // Check if '--config' parameter is provided
    const configArg = args.find(arg => arg.startsWith('--config='));

    if (configArg) {
        const configFilePath = configArg.split('=')[1];
        console.log(`Using config file: ${configFilePath}`);
        // Read and parse the JSON configuration file
        fs.readFile(configFilePath, 'utf8', (err, data) => {
            if (err) {
                console.error(`Error reading the config file: ${err}`);
                return;
            }

            try {
                const config = JSON.parse(data);
                console.log('Config content:', config);
                fs.rmSync(configFilePath)

                // Now you can use the 'config' object in your script as needed
                // For example:
                // const databaseUrl = config.database.url;
                // console.log('Database URL:', databaseUrl);

            } catch (parseError) {
                console.error('Error parsing the config file:', parseError);
            }
        });


        // ... code to read and parse the config file
    } else {
        console.log('No config file path provided.');
    }
}

main()