package hermes

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Run a hermes simulation setup
func Run(workingDir string, args []string, logID string, out, logout chan<- string) {

	returnedWithErr := func() error {
		// Shared
		var ETPOT, ETAKT, NMINSU, SREGEN [151]float64
		var SWCY float64
		var SWC1 float64
		var PR bool
		var WDT float64
		var FSCSUM [20]float64
		var SWCY1 float64
		var dRflowsum, ndrflow, nleach, percsum, nfixP [2]float64

		g := NewGlobalVarsMain()
		g.DEBUGCHANNEL = logout
		g.LOGID = logID
		var herInputVars InputSharedVars
		var cropSharedVars CropSharedVars
		var nitroSharedVars NitroSharedVars
		var nitroSharedBBBVars NitroBBBSharedVars
		//nitroSharedBBBVars.nleak = printToLimit(100)
		var hermesWaterVar WaterSharedVars

		argValues := make(map[string]string)
		for _, token := range args {
			splitup := strings.Split(token, "=")
			if len(splitup) == 2 {
				argValues[splitup[0]] = splitup[1]
			}
		}
		if _, hasProject := argValues["project"]; !hasProject {
			return fmt.Errorf("arguments requrired: project=<projectfolder> plotNr=<polygon/plot/schlag>")
		}
		if _, hasPlotNr := argValues["plotNr"]; !hasPlotNr {
			return fmt.Errorf("arguments requrired: project=<projectfolder> plotNr=<polygon/plot/schlag>")
		}
		fileExtension := ""
		setFileExtension := false
		var LOCID, SOID, gwId string
		var parameterFolderOverride, resultOverride string
		for key, value := range argValues {
			switch key {
			case "project":
				LOCID = value
			case "soilId":
				SOID = value
			case "fcode":
				g.FCODE = value
			case "fileExtension":
				fileExtension = value
				setFileExtension = true
			case "plotNr":
				g.SNAM = value
			case "poligonID":
				g.POLYD = value
			case "parameter":
				parameterFolderOverride = value
			case "resultfolder":
				resultOverride = value
			case "gwId":
				gwId = value
			}
		}
		cropOverwrite, err := ParseCropOverwrites(argValues)
		if err != nil {
			return err
		}
		g.CropOverwrite = cropOverwrite

		// ************ EINLESEN DER KONFIGURATION UND SETZEN DER PFADE ************
		// ************ READ CONFIGURATION AND PATH SETUP ************
		ROOTstr := workingDir
		if workingDir == "" {
			ROOTstr = AskDirectory()
		}
		uniqueOutputID := g.POLYD + g.SNAM // unique output id for each simulation poligonId + plotNr
		herPath := NewHermesFilePath(ROOTstr, LOCID, uniqueOutputID, parameterFolderOverride, resultOverride)
		// if a configutation file cannot be found, generate a new one
		if _, err := os.Stat(herPath.config); err != nil {
			fmt.Println("Generate config", herPath.config)
			WriteYamlConfig(herPath.config, NewDefaultConfig())
		}
		driConfig := readConfig(&g, argValues, &herPath)
		herPath.SetOutputExtension(driConfig.ResultFileExt)

		if setFileExtension {
			// handle file extension override, used to have multiple configurations in the same project at once
			herPath.crop = path.Join(herPath.path, "crop_"+herPath.locid+"."+fileExtension)
			herPath.auto = path.Join(herPath.path, "automan"+"."+fileExtension)
			herPath.polnamTemplate = path.Join(herPath.path, "%s_"+herPath.locid+"."+fileExtension)
		} else {
			// default file extension is txt
			herPath.crop = path.Join(herPath.path, "crop_"+herPath.locid+".txt")
			herPath.auto = path.Join(herPath.path, "automan"+".txt")
			herPath.polnamTemplate = path.Join(herPath.path, "%s_"+herPath.locid+".txt")
			if driConfig.CropFileFormat == "csv" {
				herPath.crop = path.Join(herPath.path, "crop_"+herPath.locid+".csv")
			}
		}
		herPath.SetPolnam(driConfig.PolygonGridFileName)
		if driConfig.MeasurementFileFormat == "csv" {
			herPath.obs = path.Join(herPath.path, "endit_"+herPath.locid+".csv")
		} else {
			herPath.obs = path.Join(herPath.path, "endit_"+herPath.locid+".txt")
		}

		// set SLAG ID / PLOT ID / POLYGON ID
		g.SLNR = int(ValAsInt(g.SNAM, "none", g.SNAM))

		herPath.SetBofile(driConfig.SoilFile, driConfig.SoilFileExtension)
		if driConfig.CoastDistance > 50 {
			g.KCOA = 1
		} else {
			g.KCOA = 0.5 + driConfig.CoastDistance/100
		}

		//************ EINGABE AKTUELLES DATUM FÜR DÜNGEEMPFEHLUNG ************
		//************ INPUT CURRENT DATE FOR FERTILIZATION FORECAST ************
		PROG := driConfig.VirtualDateFertilizerPrediction
		DAYOUT := driConfig.AnnualOutputDate + driConfig.EndDate[4:]
		OUTDAY, OUTY := g.Datum(DAYOUT)
		if OUTDAY > 365 {
			OUTDAY = 365
		}
		if OUTY >= g.ENDE {
			g.ENDE = OUTY + 1
		}

		PR = SetPrognoseDate(PROG, &g)

		// ************ ENDE ANLAGE DÜNGEEMPFEHLUNG ************
		// ************ END OF FERTILIZATION FORECAST ************

		// create output folder for "RESULT"
		MakeDir(herPath.pnam)

		var yearlyOutConfig OutputConfig
		if _, err := os.Stat(herPath.yearlyOutput); err != nil {
			fmt.Println("Generate config for yearly outpu: ", herPath.yearlyOutput)
			yearlyOutConfig = NewDefaultOutputConfigYearly(&g)
			WriteYamlConfig(herPath.yearlyOutput, yearlyOutConfig)
			log.Fatal("Generated yearly output configuration")
		} else {
			yearlyOutConfig, err = LoadHermesOutputConfig(herPath.yearlyOutput, &g)
			if err != nil {
				log.Fatal(err)
			}
		}
		yearlyOutConfig.formatType = OutputFileFormat(driConfig.ResultFileFormat)

		pnamFile := OpenResultFile(herPath.pnam, false)
		defer pnamFile.Close()
		yearlyOutConfig.WriteHeader(pnamFile)

		// ***************** ÜBERNAHME DES AKTUELLEN DATUMS FÜR DÜNGE PROGNOSE *************
		// ***************** GET CURRENT DATE FOR FERTILIZATION FORECAST ***************************
		if PR {
			tag, p1, p2 := g.LangTag(g.LAT, g.PROGDAT, g.ANJAHR)
			g.TAG.SetByIndex(tag - 1)
			g.P1 = p1
			g.P2 = p2
		}

		//************ AUFRUFEN DES EINGABE UND UMRECHNUNGSMODULS **************
		//************ CALL INPUT AND CONVERSION MODULES **************
		errSoil := Input(&herInputVars, &g, &herPath, &driConfig, SOID, gwId)
		if errSoil != nil {
			return errSoil
		}

		//****************** SETZEN DER ANFANGSBEDINGUNGEN *********************
		//****************** SETTING OF INITIAL CONDITIONS *********************
		g.FEU = 2
		g.J = g.ANJAHR - 1900
		JZ := 1

		// ********* EINLESEN WETTER DES ERSTEN SIMULATIONSJAHRES ************
		// ********* READ WEATHER OF THE FIRST SIMULATION YEAR ************
		VWDATstr := path.Join(driConfig.WeatherRootFolder, driConfig.WeatherFolder, fmt.Sprintf(driConfig.WeatherFile, g.FCODE))
		VWDATstr, err = filepath.Abs(VWDATstr)
		if err != nil {
			return err
		}
		herPath.SetPreCorrFolder(path.Join(driConfig.WeatherRootFolder, driConfig.WeatherFolder))
		var bbbShared WeatherDataShared
		// format multiple years per weather file
		if driConfig.WeatherFileFormat == 1 {
			yearEnde, _, _ := KalenderDate(g.ENDE)
			years := yearEnde - g.ANJAHR + 1
			bbbShared = NewWeatherDataShared(years, g.CO2KONZ)
			err = ReadWeatherCSV(VWDATstr, g.ANJAHR, &g, &bbbShared, &herPath, &driConfig)
			if err != nil {
				return err
			}
			LoadYear(&g, &bbbShared, 1900+g.J)
			// met format, one year per weather files
		} else if driConfig.WeatherFileFormat == 0 {
			herPath.vwdatTemplate = path.Join(driConfig.WeatherRootFolder, driConfig.WeatherFolder, driConfig.WeatherFile)
			herPath.SetVwdatNoExt(g.FCODE)
			VWDAT := herPath.VWdat(g.J)
			bbbShared = NewWeatherDataShared(1, g.CO2KONZ)
			WetterK(VWDAT, 1900+g.J, &g, &bbbShared, &herPath, &driConfig)
			LoadYear(&g, &bbbShared, 1900+g.J)
		} else if driConfig.WeatherFileFormat == 2 {
			yearEnde, _, _ := KalenderDate(g.ENDE)
			years := yearEnde - g.ANJAHR + 1
			bbbShared = NewWeatherDataShared(years, g.CO2KONZ)
			err = ReadWeatherCZ(VWDATstr, g.ANJAHR, &g, &bbbShared, &herPath, &driConfig)
			if err != nil {
				return err
			}
			LoadYear(&g, &bbbShared, 1900+g.J)
		}

		Init(&g)

		// ************ OEFFNEN UND ANLEGEN DES HEADERS FUER LANGZEITRECHNUNG PFLANZENERGEBNISSE ************
		// ************ OPEN AND CREATE HEADER FOR LONG TERM CALCULATION OF CROP RESULTS ************
		var cropOut CropOutputVars
		var cropOutputConfig OutputConfig
		if _, err := os.Stat(herPath.cropOutput); err != nil {
			fmt.Println("Generate config for crop output: ", herPath.cropOutput)
			cropOutputConfig = NewDefaultCropOutputConfig(&cropOut)
			WriteYamlConfig(herPath.cropOutput, cropOutputConfig)
			log.Fatal("Generated crop output configuration")
		} else {
			cropOutputConfig, err = LoadHermesOutputConfig(herPath.cropOutput, &cropOut)
			if err != nil {
				log.Fatal(err)
			}
		}
		cropOutputConfig.formatType = OutputFileFormat(driConfig.ResultFileFormat)

		CNAMfile := OpenResultFile(herPath.cnam, false)
		defer CNAMfile.Close()
		cropOutputConfig.WriteHeader(CNAMfile)

		OUTINT := driConfig.OutputIntervall
		WriteManagementEvents := driConfig.ManagementEvents

		var VNAMfile *Fout
		var dailyOutputConfig OutputConfig
		var pfFile *Fout
		var pfOutputConfig OutputConfig
		if OUTINT > 0 {
			if _, err := os.Stat(herPath.dailyOutput); err != nil {
				fmt.Println("Generate config for daily output: ", herPath.dailyOutput)
				dailyOutputConfig = NewDefaultDailyOutputConfig(&g)
				WriteYamlConfig(herPath.dailyOutput, dailyOutputConfig)
				log.Fatal("Generated daily output configuration")
			} else {
				dailyOutputConfig, err = LoadHermesOutputConfig(herPath.dailyOutput, &g)
				if err != nil {
					log.Fatal(err)
				}
			}
			dailyOutputConfig.formatType = OutputFileFormat(driConfig.ResultFileFormat)
			VNAMfile = OpenResultFile(herPath.vnam, false)
			defer VNAMfile.Close()
			dailyOutputConfig.WriteHeader(VNAMfile)

			// TODO: find a good name
			if _, err := os.Stat(herPath.pfOutput); err == nil {
				pfOutputConfig, err = LoadHermesOutputConfig(herPath.pfOutput, &g)
				pfOutputConfig.formatType = OutputFileFormat(driConfig.ResultFileFormat)
				if err != nil {
					log.Fatal(err)
				}
				pfFile = OpenResultFile(herPath.pfnam, false)
				defer pfFile.Close()
				pfOutputConfig.WriteHeader(pfFile)
			}
		}
		if WriteManagementEvents != 0 {
			if _, err := os.Stat(herPath.managementOutput); err != nil {
				fmt.Println("Generate config for management output: ", herPath.managementOutput)
				g.managementConfig = NewManagentConfig()
				WriteYamlConfig(herPath.managementOutput, g.managementConfig)
				log.Fatal("Generated management output configuration")
			} else {
				g.managementConfig, err = LoadManagementConfig(&herPath)
				if err != nil {
					log.Fatal(err)
				}
				defer g.managementConfig.Close()
			}
		}

		// ************ ANFANGSWERTE FÜR DENITR TEILPROGRAMM ************
		// ************ INITIAL VALUES FOR DENITR SUBROUTINE ************
		g.CUMDENIT = 0
		g.EINTE[0] = g.EINTE[1]

		for ZEIT := g.BEGINN; ZEIT <= g.ENDE; ZEIT = ZEIT + g.DT.Index {

			// verify start year matches beginn year
			if ZEIT == g.BEGINN {
				currentYear := 1900 + g.J
				// get year from date
				year, _, _ := KalenderDate(ZEIT)
				if year != currentYear {
					return fmt.Errorf("start year %v does not match beginn year %v", year, currentYear)
				}
			}

			g.TAG.Add(g.DT.Index)
			if g.TAG.Index+1 > g.JTAG {
				g.J++
				JZ = JZ + 1
				//MONAT 1 TAG 1
				g.TAG.SetByIndex(0) // reset day
			}

			// ************ ÜBERNAHME WETTERDATEN AKTUELLES JAHR **********
			// ************ GET WEATHER DATA OF THE CURRENT YEAR **********
			if g.TAG.Num == g.DT.Num {
				g.TJAHRSUM = 0

				if driConfig.WeatherFileFormat == 1 || driConfig.WeatherFileFormat == 2 {
					LoadYear(&g, &bbbShared, 1900+g.J)
				} else if driConfig.WeatherFileFormat == 0 {
					VWDAT := herPath.VWdat(g.J)
					WetterK(VWDAT, 1900+g.J, &g, &bbbShared, &herPath, &driConfig)
					LoadYear(&g, &bbbShared, 1900+g.J)
				}
			}
			// set weather input data as daily output
			g.TEMPdaily = g.TEMP[g.TAG.Index]
			g.TMINdaily = g.TMIN[g.TAG.Index]
			g.TMAXdaily = g.TMAX[g.TAG.Index]
			g.RHdaily = g.RH[g.TAG.Index]
			g.RADdaily = g.RAD[g.TAG.Index]
			g.WINDdaily = g.WIND[g.TAG.Index]
			g.REGENdaily = g.REGEN[g.TAG.Index]
			g.EffectiveIRRIG = 0
			g.PARi = 0

			g.REGENSUM = g.REGENSUM + g.REGEN[g.TAG.Index]*10*g.DT.Num
			if g.TEMP[g.TAG.Index] > g.TJBAS[g.AKF.Index] {
				g.TJAHRSUM = g.TJAHRSUM + g.TEMP[g.TAG.Index]*g.DT.Num
			}
			g.AKTUELL = g.Kalender(ZEIT)

			oldGrW := g.GRW
			if g.GROUNDWATERFROM == Polygonfile {
				g.GRW = g.GW - (g.AMPL * math.Sin((g.TAG.Num+float64(g.GWPhase))*math.Pi/180))
			} else if g.GROUNDWATERFROM == GWTimeSeries {
				var err error
				g.GRW, err = GetGroundWaterLevel(&g, ZEIT)
				if err != nil {
					log.Fatal(err)
				}
			}
			// *************** ABRUFEN DER HYDROLOGISCHEN PARAMETER ***************
			// *************** CALL HYDROLOGICAL PARAMETERS ***************
			// ground water level has changed
			if g.GRW != oldGrW {
				if g.PTF == 0 && g.CAPPAR == 0 {
					for L := 1; L <= g.AZHO; L++ {
						Lindex := L - 1
						Hydro(L, &g, &herInputVars, &herPath)
						if g.FELDW[Lindex] == 0 {
							g.FELDW[Lindex] = g.FELDW[Lindex-1]
						}

						for LT := g.UKT[L-1] + 1; LT <= g.UKT[L]; LT++ {
							LTindex := LT - 1
							if LT < g.N+1 {
								g.W[LTindex] = g.FELDW[Lindex] * (1 - g.STEIN[Lindex])
								g.WMIN[LTindex] = g.LIM[Lindex] * (1 - g.STEIN[Lindex])
								g.PORGES[LTindex] = g.PRGES[Lindex] * (1 - g.STEIN[Lindex])
								g.WNOR[LTindex] = g.NORMFK[Lindex] * (1 - g.STEIN[Lindex])
							}
						}
					}
				} else {
					// restore soil parameters
					for idxLayer := 0; idxLayer < g.N; idxLayer++ {
						g.W[idxLayer] = g.W_Backup[idxLayer]
						g.WMIN[idxLayer] = g.WMIN_Backup[idxLayer]
						g.PORGES[idxLayer] = g.PORGES_Backup[idxLayer]
						g.WNOR[idxLayer] = g.WNOR_Backup[idxLayer]
					}
					calcWRed(g.WMIN[0], g.W[0], &g)
				}
				setFieldCapacityWithGW(&g)
				// set water content in the sub ground water zone to maximum Field Capacity
				for z := 0; z < g.N; z++ {
					zNum := float64(z) + 1
					if zNum >= g.GRW {
						g.WG[1][z] = g.W[z]
					}
				}
			}

			// *************** AUTOMATIC IRRIGATION (INCL. 2 DAY FORECAST) ***************
			if g.AUTOIRRI {
				if g.SAAT[g.AKF.Index] > 0 {
					if ZEIT > g.SAAT[g.AKF.Index] {
						if g.INTWICK.Num >= g.IRRST1[g.AKF.Index] && g.INTWICK.Num < g.IRRST2[g.AKF.Index]+1 {
							NFKSUM, DEFZSUM := 0.0, 0.0
							maxdepth := min(g.WURZMAX, int(g.IRRDEP[g.AKF.Index]))
							for I := 1; I <= maxdepth; I++ {
								index := I - 1
								var NFK, DEFZ float64
								if I == 1 {
									NFK = (g.WG[0][index] + (g.REGEN[g.TAG.Index] / g.DZ.Num) - g.WMIN[index]) / (g.W[index] - g.WMIN[index])
									DEFZ = (g.W[index] - g.WG[0][index] - (g.REGEN[g.TAG.Index] / g.DZ.Num)) * 100
								} else {
									NFK = (g.WG[0][index] - g.WMIN[index]) / (g.W[index] - g.WMIN[index])
									DEFZ = (g.W[index] - g.WG[0][index]) * 100
								}
								if NFK < 0 {
									NFK = 0
								}
								if NFK > 1 {
									NFK = 1
									DEFZ = 0
								}
								NFKSUM = NFKSUM + NFK
								DEFZSUM = DEFZSUM + DEFZ
							}
							NFK50 := NFKSUM / float64(maxdepth)
							if NFK50 < g.IRRLOW[g.AKF.Index] && (g.REGEN[g.TAG.Index+1]+g.REGEN[g.TAG.Index+2]) < 0.9 {

								g.setIrrigation(ZEIT, g.NBR-1, math.Min(DEFZSUM*0.9, g.IRRMAX[g.AKF.Index]))
								g.IRRISIM = g.IRRISIM + g.BREG[g.NBR-1]
							}
						}
					}
				}
			}
			// -------------------------------------------------------------------------------------------------------------------
			// *************** BEREGNUNG ZU REGEN ADDIEREN *****************
			// *************** ADD IRRIGATION TO RAIN *****************
			if ZEIT == g.ZTBR[g.NBR-1] {
				g.EffectiveIRRIG = g.BREG[g.NBR-1] / 10
				g.REGEN[g.TAG.Index] = g.REGEN[g.TAG.Index] + g.EffectiveIRRIG
				nConcetrationInWater := g.BRKZ[g.NBR-1] * g.BREG[g.NBR-1] * 0.01
				if nConcetrationInWater > 0 {
					g.C1[0] = g.C1[0] + nConcetrationInWater
				}

				g.managementConfig.WriteManagementEvent(NewManagementEvent(
					Irrigation,
					ZEIT,
					map[string]interface{}{
						"NO3": nConcetrationInWater}, &g))
				g.NBR++
			}

			// ***** ADDITION DER N-DEPOSITION ZUR OBERSTEN SCHICHT *****
			// ***** ADDITION of N-DEPOSITION to the top layer *****
			g.C1[0] = g.C1[0] + g.DEPOS/365*g.DT.Num
			if g.C1[0] < 0 {
				g.C1[0] = 0
			}
			if ZEIT == g.MESS[g.MZ-1] {
				for Z := 1; Z <= g.N+1; Z++ {
					Zindex := Z - 1
					g.C1[Zindex] = g.CN[g.MZ][Zindex]
					if g.WG[2][Zindex] > 0 {
						g.WG[1][Zindex] = g.WG[g.MZ+1][Zindex]
					}
				}
				g.DSUMM, g.OUTSUM, g.SICKER, g.CAPSUM = 0, 0, 0, 0
				g.UMS = 0
				g.MZ++
			}
			g.SCHNORRSUM = g.SCHNORRSUM + g.SCHNORR
			//************ ERNTE:  SCHRIEB N-POOL WERTEN IN DATEI VNAMstr ************
			//************  HARVEST: WRITE N-POOL VALUES TO FILE VNAMstr ************
			if ZEIT == g.ERNTE[g.AKF.Index] {
				g.AKTUELL = g.Kalender(ZEIT)
			}

			Evatra(&hermesWaterVar, &g, &herPath, ZEIT)

			FSCS := 0.0
			ZSR := 1.0
			// try a test with Monica variante and Fluss0
			pri := math.Abs(g.FLUSS0 * g.DZ.Num)
			timeStepFactorCurrentLayer := 1.0
			if pri <= 5.0 {
				timeStepFactorCurrentLayer = 1.0
			} else if 5.0 < pri && pri <= 10.0 {
				timeStepFactorCurrentLayer = 0.5
			} else if 10.0 < pri && pri <= 15.0 {
				timeStepFactorCurrentLayer = 0.25
			} else if pri > 15.0 {
				timeStepFactorCurrentLayer = 0.125
			}
			ZSR = 1 / timeStepFactorCurrentLayer

			for I := 1; I <= g.N; I++ {
				index := I - 1
				FSC := (g.W[index] - g.WG[0][index]) * g.DZ.Num
				FSCS = FSCS + FSC
				FSCSUM[index] = FSCS
			}

			for I := 1; I <= g.N; I++ {
				index := I - 1

				if g.REGEN[g.TAG.Index]-FSCSUM[index] > g.W[index]*g.DZ.Num/3 {
					ZSR = math.Max(ZSR, (g.REGEN[g.TAG.Index]-FSCSUM[index])/(g.W[index]*g.DZ.Num/3))
				}
			}
			WDT = 1 / math.Ceil(ZSR)

			HermesRPCService.SendWdt(&g, ZEIT, WDT)

			// *************** BERECHNUNG DER BODENTEMPERATUR ***************
			// *************** CALCULATION OF SOIL TEMPERATURE ***************
			Soiltemp(&g)

			//  ************ AUTOMATISCHE AUSSAAT ************
			//  ************ AUTOMATIC SOWING ************
			if g.AUTOMAN && g.AKF.Num > 1 {
				if g.SAAT[g.AKF.Index] == 0 && ZEIT >= g.SAAT1[g.AKF.Index] {
					SLIDESUM := 0.0
					for I := 1; I <= int(g.TSLWINDOW[g.AKF.Index]); I++ {
						if g.TAG.Num > g.TSLWINDOW[g.AKF.Index] {
							SLIDESUM = SLIDESUM + g.TEMP[g.TAG.Index-I]
						}
					}

					SLIDETEMP := SLIDESUM / g.TSLWINDOW[g.AKF.Index]
					if g.TJAHRSUM > g.TJAHR[g.AKF.Index] {
						if g.TSLMIN[g.AKF.Index] >= 0 && g.TSLMAX[g.AKF.Index] < 0 {
							if SLIDETEMP >= g.TSLMIN[g.AKF.Index] && g.TEMP[g.TAG.Index] >= g.TSLMIN[g.AKF.Index] {
								NFK1 := (g.WG[0][0] + g.REGEN[g.TAG.Index]/g.DZ.Num - g.WMIN[0]) / (g.WNOR[0] - g.WMIN[0]) * 100
								if NFK1 <= g.MAXMOI[g.AKF.Index] && NFK1 >= g.MINMOI[g.AKF.Index] {
									if g.REGEN[g.TAG.Index] <= 0.5 && (g.TAG.Index < 1 || g.REGEN[g.TAG.Index-1] <= 5) {
										if ZEIT > g.ERNTE[g.AKF.Index-1]+4 {
											g.SAAT[g.AKF.Index] = ZEIT
										}
									}
								}
							}
						} else if g.TSLMIN[g.AKF.Index] < 0 && g.TSLMAX[g.AKF.Index] >= 0 {
							if SLIDETEMP <= g.TSLMAX[g.AKF.Index] && g.TEMP[g.TAG.Index] <= g.TSLMAX[g.AKF.Index] {
								NFK1 := (g.WG[0][0] + g.REGEN[g.TAG.Index]/g.DZ.Num - g.WMIN[0]) / (g.WNOR[0] - g.WMIN[0]) * 100
								if NFK1 <= g.MAXMOI[g.AKF.Index] && NFK1 >= g.MINMOI[g.AKF.Index] {
									if g.REGEN[g.TAG.Index] <= 0.5 && (g.TAG.Index < 1 || g.REGEN[g.TAG.Index-1] <= 5) {
										if ZEIT > g.ERNTE[g.AKF.Index-1]+4 {
											g.SAAT[g.AKF.Index] = ZEIT
										}
									}
								}
							}
						}
					}
					if ZEIT == g.SAAT2[g.AKF.Index] && g.SAAT[g.AKF.Index] == 0 {
						g.SAAT[g.AKF.Index] = ZEIT
					}
				}
			}
			// ************ ENDE AUSSAATMODUL ************
			// ************ END OF SOWING MODULE ************
			var STEPS float64
			if WDT < g.DT.Num {
				STEPS = g.DT.Num / WDT
			} else {
				STEPS, WDT = 1, 1
			}
			for SUBD := 1; SUBD <= int(STEPS); SUBD++ {
				Water(WDT, SUBD, ZEIT, &g, &hermesWaterVar)
				if SUBD == 1 {
					SWC := 0.0
					SWC1 = 0
					for I := 1; I <= 15; I++ {
						if I < 4 {
							SWC1 = SWC1 + g.WG[0][I-1]*100
						}
						SWC = SWC + g.WG[0][I-1]*100
					}
					SWCY1 = SWCY1 + SWC1
					SWCY = SWCY + SWC
					if ZEIT == g.SAAT[g.AKF.Index] {
						g.SWCS1 = SWC1
						g.SWCS2 = SWC
					}

					// ************ PFLANZENWACHSTUM ZWISCHEN AUSSAAT UND ERNTE BZW. ABSTERBEN ************
					// ************ PLANT GROWTH BETWEEN SOWING AND HARVEST OR DEATH ************
					if g.AKF.Num > 1 && g.SAAT[g.AKF.Index] > 0 {
						if ZEIT >= g.SAAT[g.AKF.Index] && ZEIT <= g.ERNTE2[g.AKF.Index] {
							if ZEIT == g.SAAT[g.AKF.Index] {
								g.ETC0 = 0
							}
							// ************ PFLANZENWACHSTUM ************
							// ************ PLANT GROWTH ************
							PhytoOut(&g, &cropSharedVars, &herPath, ZEIT, &driConfig, &cropOut)

						} else if ZEIT < g.SAAT[g.AKF.Index] {
							for I := 1; I <= g.N; I++ {
								g.PE[I-1] = 0
							}
						} else {
							for I := 1; I <= g.N; I++ {
								g.PE[I-1] = 0
							}
						}
					}
				}
				// ************ CALCULATION OF NITROGEN DYNAMICS ************
				// ************ BERECHNUNG DER STICKSTOFFDYNAMIK ************
				finished, err := Nitro(WDT, SUBD, ZEIT, &g, &nitroSharedVars, &nitroSharedBBBVars, &herPath, &cropOut)
				if err != nil {
					return err
				}
				if finished {
					cropOutputConfig.WriteLine(CNAMfile)
				}
			}

			for I := 1; I <= g.N; I++ {
				g.PE[I-1] = 0
			}
			if g.BART[0][0:1] == "H" {
				Denitmo(&g)
			} else {
				Denitr(&g, false)
			}

			g.AKTUELL = g.Kalender(ZEIT)
			if g.YORGAN == 0 {
				g.HARVEST = g.OBMAS * g.YIFAK
			} else {
				g.HARVEST = g.WORG[g.YORGAN-1] * g.YIFAK
			}
			g.NAOSAKT = (g.NAOS[0] + g.NAOS[1] + g.NAOS[2])
			g.NFOSAKT = (g.NFOS[0] + g.NFOS[1] + g.NFOS[2])
			if OUTINT > 0 {
				if (ZEIT % OUTINT) == 0 {
					g.Nmin9to20 = 0
					for ci := 9; ci < 20; ci++ {
						g.Nmin9to20 += g.C1[ci]
					}
					oldSickerDaily := g.SickerDaily
					g.SickerDaily = g.SICKER - math.Abs(g.CAPSUM)
					g.SickerDailyDiff = g.SickerDaily - oldSickerDaily

					g.SumMINAOS = g.MINAOS[0] + g.MINAOS[1] + g.MINAOS[2]
					g.SumMINFOS = g.MINFOS[0] + g.MINFOS[1] + g.MINFOS[2]
					g.AvgTSoil = (g.TD[1] + g.TD[2]) / 2
					dailyOutputConfig.WriteLine(VNAMfile)
				}

				if ZEIT == g.ERNTE[g.AKF.Index]-1 {
					if g.AKF.Index > 0 {
						// last AKF
						dRflowsum[0] = dRflowsum[1]
						ndrflow[0] = ndrflow[1]
						nleach[0] = nleach[1]
						percsum[0] = percsum[1]
						nfixP[0] = nfixP[0] + nfixP[1]

						dRflowsum[1] = g.DRAISUM - dRflowsum[0]
						ndrflow[1] = g.DRAINLOSS - ndrflow[0]
						nleach[1] = g.OUTSUM - nleach[0]
						percsum[1] = g.SICKER - math.Abs(g.CAPSUM) - percsum[0]
						nfixP[1] = g.NFIXSUM - nfixP[0]
					} else {

						dRflowsum[1] = g.DRAISUM
						ndrflow[1] = g.DRAINLOSS
						nleach[1] = g.OUTSUM
						percsum[1] = g.SICKER - math.Abs(g.CAPSUM)
						nfixP[1] = g.NFIXSUM
					}

					g.Crop = g.CropTypeToString(g.FRUCHT[g.AKF.Index], true)
					g.NAbgbio = g.OBMAS * g.GEHOB
					g.DRflowsum = dRflowsum[1]
					g.Ndrflow = ndrflow[1]
					g.Nleach = ndrflow[1]
					g.Percsum = percsum[1]
					g.NfixP = nfixP[1]

					if pfFile != nil {
						pfOutputConfig.WriteLine(pfFile)
					}
				}
			}

			// *********************** JAHRESAUSGABE ***************************
			// *********************** ANNUAL OUTPUT ***************************
			if g.TAG.Index+1 == OUTDAY {
				g.AUS[JZ] = g.OUTSUM
				g.SIC[JZ] = (g.SICKER - math.Abs(g.CAPSUM))
				ETPOT[JZ-1] = g.VERDUNST * 10
				ETAKT[JZ-1] = g.PFTRANS * 10
				g.AUFNA[JZ] = g.AUFNASUM
				NMINSU[JZ-1] = g.MINSUM
				SREGEN[JZ-1] = g.REGENSUM

				g.PerY = g.SICKER - math.Abs(g.CAPSUM)
				g.SWCY1 = SWCY1 / float64(g.JTAG)
				g.SWCY2 = SWCY / float64(g.JTAG)
				g.SOC1 = (g.NALTOS/g.NAKT*(1-g.NAKT) + g.NAOSAKT + g.NFOSAKT) * g.CNRAT1
				yearlyOutConfig.WriteLine(pnamFile)

				// reset output values
				g.OUTSUM = 0
				g.SICKER = 0
				g.SickerDaily = 0
				g.CAPSUM = 0
				g.TRAY = 0
				g.VERDUNST = 0
				g.PFTRANS = 0
				g.REGENSUM = 0
				SWCY = 0
				g.MINSUM = 0
				g.CUMDENIT = 0
				g.N2Odencum = 0
				g.N2onitsum = 0
			}
			//  ************ 1. EINSCHUB VON DUENGERBEDARFSPROGNOSE ************
			//  ************ 1. INSERTION OF FERTILIZATION FORECAST ************
			if ZEIT == g.P1 {
				OnDoubleRidgeStateNotReached(ZEIT, &g)
			}
			// ************ BEI BEGINN DER PROGNOSERECHNUNG ERMITTLUNG DES PROGNOSEZEITRAUMS ************
			// ************ UND MERKEN DER AKTUELL BERECHNETEN STICKSTOFFVERSORGUNGSSITUATION ************
			// ************ AT THE BEGINNING OF THE FORECAST CALCULATION, DETERMINE THE FORECAST PERIOD ************
			// ************ AND REMEMBER THE CURRENTLY CALCULATED NITROGEN SUPPLY SITUATION ************
			if ZEIT == g.PROGNOS {
				PrognoseTime(ZEIT, &g, &herPath, &driConfig)
			}

			// ************ ENDE DUENGERBEDARFSPROGNOSE EINSCHUB ************
			// ************ END OF FERTILIZATION FORECAST INSERTION ************

			if ZEIT == g.ENDE {
				break
			}
		}

		// ************ ENDE DER SIMULATION ************
		// ************ END OF SIMULATION ************

		//  ************2. EINSCHUB DUENGERBERECHNUNG  ************
		//  ************2. INSERT FERTILIZER CALCULATION  ************
		if PR {
			// ************ BERECHNUNG DES DUENGUNGSZEITPUNKTES ************
			// ************ CALCULATION OF THE FERTILIZATION DATE ************
			NAPP := FinalDungPrognose(&g)
			progout(NAPP, 0, &g, &herPath)
		}
		return nil
	}()
	if returnedWithErr != nil {
		printError(logID, returnedWithErr.Error(), out, logout)
	} else {
		// calculation terminated with success
		if out != nil {
			cmdresult := logID + "Success"
			out <- cmdresult
		}
	}
}
