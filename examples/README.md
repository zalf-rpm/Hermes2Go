# To run Hermes

you run Hermes in a terminal

1. Run the project as definded (read from project/modinp.txt):  
hermes2go.exe

2. Run Hermes2Go in batch mode:  
hermes2go.exe -module batch -logoutput -concurrent 1 -batch myP_batch.txt

Batch mode allows to run multiple project setups concurrently.  
-batch myP_batch.txt - defines the batch file to be executed.  
-concurrent 1 - defines the possible number of parallel executions. A good value is the number of processors of your computer.  
-logoutput - will write debug output to the console. If you omit this parameter no output will be written to the console.  

The batch file contains one line for each project to be executed, e.g:

project=myP WeatherFolder=historical soilId=075 plotNr=10001 Altitude=73 Latitude=52.6732 poligonID=29872 resultfolder=RESULT_myP  
project=myP WeatherFolder=historical soilId=075 plotNr=10002 Altitude=73 Latitude=52.6732 poligonID=29872 resultfolder=RESULT_myP  

A batch file allows you to change some attributes of the project, e.g. the plotNr, weather data etc.
The project name must be the same as the name of the project folder.
For example, myP is the name of a sample project, which is the same as the name of the folder containing the project files.
The project folder contains the following files:

- project\myP\automan.txt - contains values for automatic management
- project\myP\config.yml - contains the configuration of the project
- project\myP\crop_myP.txt - contains the crop rotation
- project\myP\cropout_conf.yml - contains the configuration of the crop output
- project\myP\dailyout_conf.yml - contains the configuration of the daily output
- project\myP\endit_myP.txt - contains initial values for the simulation
- project\myP\fert_myP.txt - contains the fertilization plan
- project\myP\irr_myP.txt - contains the irrigation plan
- project\myP\poly_myP.txt - contains the polygon/Field mapping data
- project\myP\soil_myP.csv - contains the soil data
- project\myP\til_myP.txt - contains the tillage plan
- project\myP\yearlyout_conf.yml - contains the configuration of the yearly output

The output files are stored in the folder RESULT_myP. The output files are for example:  
RESULT_myP\C2987210001.csv  
RESULT_myP\V2987210001.csv  
RESULT_myP\Y2987210001.csv  

starting with the letter C, V or Y, followed by the polygon ID and the plot number.  
C2987210001.csv contains the crop output for polygon 29872 and plot 10001.  
V2987210001.csv contains the daily output for polygon 29872 and plot 10001.  
Y2987210001.csv contains the yearly output for polygon 29872 and plot 10001.  
