To run Hermes executable on windows, you need to open a terminal:

To execute the project as definded (read from project/modinp.txt), run the following command:
.\hermes2go.exe

To execute in batch mode, run the following command:
.\hermes2go.exe -module batch -logoutput -concurrent 1 -batch myP_batch.txt


The batch mode allows to run multiple projects in parallel. The number of parallel projects is defined by the parameter -concurrent 1.
The parameter -logoutput will write debug output to the console. If you omit this parameter no output will be written to the console.
The parameter -batch myP_batch.txt defines the batch file to be executed.

The batch file contains one line for each project to be executed. Each line contains the following parameters:
project=myP WeatherFolder=historical soilId=075 plotNr=10001 Altitude=73 Latitude=52.6732 poligonID=29872 resultfolder=RESULT_myP
project=myP WeatherFolder=historical soilId=075 plotNr=10002 Altitude=73 Latitude=52.6732 poligonID=29872 resultfolder=RESULT_myP

A batch file allows you to change some attribute of the project, e.g. the plotNr, weather data etc. The project name must be the same as the name of the project folder.

myP is the name of the project, which is the same as the name of the folder containing the project files. 
The project folder contains the following files:

project\myP\automan.txt - contains values for automatic management
project\myP\config.yml - contains the configuration of the project
project\myP\crop_myP.txt - contains the crop rotation
project\myP\cropout_conf.yml - contains the configuration of the crop output
project\myP\dailyout_conf.yml - contains the configuration of the daily output
project\myP\endit_myP.txt - contains initial values for the simulation
project\myP\fert_myP.txt - contains the fertilization plan
project\myP\irr_myP.txt - contains the irrigation plan
project\myP\poly_myP.txt - contains the polygon/Field mapping data
project\myP\soil_myP.csv - contains the soil data
project\myP\til_myP.txt - contains the tillage plan
project\myP\yearlyout_conf.yml - contains the configuration of the yearly output

The output files are stored in the folder RESULT_myP. The output files are:
RESULT_myP\C2987210001.csv
RESULT_myP\C2987210002.csv
RESULT_myP\V2987210001.csv
RESULT_myP\V2987210002.csv
RESULT_myP\Y2987210001.csv
RESULT_myP\Y2987210002.csv

starting with the letter C, V or Y, followed by the polygon ID and the plot number.
C2987210001.csv contains the crop output for polygon 29872 and plot 10001.
V2987210001.csv contains the daily output for polygon 29872 and plot 10001.
Y2987210001.csv contains the yearly output for polygon 29872 and plot 10001.
