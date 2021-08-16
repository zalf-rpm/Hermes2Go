package hermes

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

func progout(NAPP int, ABZEIT int, g *GlobalVarsMain, hPath *HFilePath) {
	//CALL TC_Init
	//CALL TC_Getscreensize(ls,rs,bs,ts)
	//CALL TC_Win_Create (win_23, "CLOSE|SIZE|TITLE", .25, .75, .25, .75)
	//CALL TC_Win_SetTitle (win_23, "Empfehlung HERMES")
	//CALL TC_Show (win_23)
	//SET COLOR "white"                  !22
	//BOX AREA 0,1,0,1
	//CALL TC_Win_Active (win_23)
	//CALL TC_PushBtn_Create (pbid_7, "schließen", .6, .9, .05, .1)

	// LET DUNG1,DUNG2 = 0
	var DUNG1, DUNG2 float64
	// CLEAR
	// SET COLOR "blue"

	// IF ABZEIT > 0 THEN
	if ABZEIT > 0 {
		//    LET ENDE = ABZEIT
		g.ENDE = ABZEIT
		// END IF
	}
	// CALL KALENDER(INT(ENDE),ENDDAT$)
	ENDDAT := g.Kalender(g.ENDE)
	var SUMN1, SUMN2, SUMN3, SUMN4 float64
	var SUMA1, SUMA2, SUMA3, SUMA4 float64
	var WASS1, WASS2, WASS3 float64
	var NAPPDAT2 string
	// FOR Z= 1 TO N
	for z := 0; z < g.N; z++ {
		// IF Z < 4 THEN
		if z < 3 {
			//LET SUMN1 = SUMN1 + C1(Z)
			SUMN1 = SUMN1 + g.C1[z]
			// LET SUMA1 = SUMA1 + CA(Z)
			SUMA1 = SUMA1 + g.CA[z]
			// LET WASS1 = WASS1 + WG(1,Z)
			WASS1 = WASS1 + g.WG[1][z]
			// ELSE IF Z < 7 THEN
		} else if z < 6 {
			// LET SUMN2 = SUMN2 + C1(Z)
			SUMN2 = SUMN2 + g.C1[z]
			// LET SUMA2 = SUMA2 + CA(Z)
			SUMA2 = SUMA2 + g.CA[z]
			// LET WASS2 = WASS2 + WG(1,Z)
			WASS2 = WASS2 + g.WG[1][z]
			// ELSE IF Z < 10 THEN
		} else if z < 9 {
			// LET SUMN3 = SUMN3 + C1(Z)
			SUMN3 = SUMN3 + g.C1[z]
			// LET SUMA3 = SUMA3 + CA(Z)
			SUMA3 = SUMA3 + g.CA[z]
			// LET WASS3 = WASS3 + WG(1,Z)
			WASS3 = WASS3 + g.WG[1][z]
			// ELSE
		} else {
			// LET SUMN4 = SUMN4 + C1(Z)
			SUMN4 = SUMN4 + g.C1[z]
			// LET SUMA4 = SUMA4 + CA(Z)
			SUMA4 = SUMA4 + g.CA[z]
			// END IF
		}
		// NEXT Z
	}
	// // LET WASS1 = WASS1/3
	// WASS1 = WASS1 / 3
	// // LET WASS2 = WASS2/3
	// WASS2 = WASS2 / 3
	// // LET WASS3 = WASS3/3
	// WASS3 = WASS3 / 3
	// let prognodat$ = Progdat$(1:2) & "." & progdat$(3:4) & "." & progdat$(5:6)
	//prognodat := g.PROGDAT[0:2] + "." + g.PROGDAT[2:4] + "." + g.PROGDAT[4:]
	_, masdat := g.Datum(g.PROGDAT)
	progDatStr := g.Kalender(masdat)
	// print Endstadium$
	fmt.Println(g.ENDSTADIUM)
	// !PRINT "        ----------------------------------------------------------------"
	// PRINT "        ----------------------------------------------------------------"
	fmt.Println("        ----------------------------------------------------------------")
	// PRINT USING "        |    N Bedarfsberechnung vom  ########## von Flaeche <#####    |":PROGnoDAT$,slnr         !lnam$
	fmt.Printf("        |    N Bedarfsberechnung vom  %s von Flaeche > %05d    |\n", progDatStr, g.SLNR)
	// PRINT USING "        |            calculation of N demand  ######## of #########    |":PROGDAT$,slnam$
	fmt.Printf("        |            calculation of N demand  %s of %8s    |\n", progDatStr, g.SLNAM)
	// PRINT "        ----------------------------------------------------------------"
	fmt.Println("        ----------------------------------------------------------------")
	// !PRINT
	// PRINT" "
	fmt.Println(" ")
	// !PRINT" für Prognoserechnung wird das langjährige Mittel des Standortes verwendet"
	// !PRINT
	// IF PROGNOS < ENDE THEN
	if g.PROGNOS < g.ENDE {
		// PRINT" aktuelle Nmin-Verteilung zum ";PROGnoDAT$;":"
		// ! PRINT" simulated Nmin distribution by ";PROGDAT$;":"
		// PRINT USING "                                             0 _ 30 cm: ### kg N/ha":SUMA1
		fmt.Printf("                                             0 _ 30 cm: %03d kg N/ha\n", int(SUMA1))
		// PRINT USING "                                            30 _ 60 cm: ### kg N/ha":SUMA2
		fmt.Printf("                                            30 _ 60 cm: %03d kg N/ha\n", int(SUMA2))
		// PRINT USING "                                            60 _ 90 cm: ### kg N/ha":SUMA3
		fmt.Printf("                                            60 _ 90 cm: %03d kg N/ha\n", int(SUMA3))
		// PRINT "                                            -----------------------"
		fmt.Println("                                            -----------------------")
		// PRINT USING "                                             0 _ 90 cm: ### kg N/ha":SUMA1+SUMA2+SUMA3
		fmt.Printf("                                             0 _ 90 cm: %03d kg N/ha\n", int(SUMA1+SUMA2+SUMA3))
		// PRINT
		fmt.Println(" ")
		// PRINT USING "Es wurden bereits ### kg N/ha durch die Pflanzen aufgenommen":PLANA
		fmt.Printf("Es wurden bereits %03d kg N/ha durch die Pflanzen aufgenommen\n", int(g.PLANA))
		// !PRINT USING "there are already ### kg N/ha taken up by crops":PLANA
		//!IF ENDE < ENDPRO THEN
		//IF ABZEIT <> 0 THEN
		if ABZEIT != 0 {
			//PRINT "*************************************************************************"
			fmt.Println("*************************************************************************")
			//PRINT       "|   Prognose des Stickstoffbedarfs bis zum Abbruchtermin                 |"
			fmt.Println("|   Prognose des Stickstoffbedarfs bis zum Abbruchtermin                 |")
			//PRINT USING "|   Das N-Defizit bis zum ########   betraegt :                          |":ENDDAT$
			fmt.Printf("|   Das N-Defizit bis zum %s   betraegt :                          |\n", ENDDAT)
			// !PRINT       "|   Prediction of nitrogen deficit until user break                     |"
			// !PRINT USING "|   the N deficit until ##########   will be :                          |":ENDDAT$
			//IF DUNGBED < 2 THEN
			if g.DUNGBED < 2 {
				//PRINT "|   Vorrat bislang ausreichend, keine Duengung notwendig                 |"
				fmt.Println("|   Vorrat bislang ausreichend, keine Duengung notwendig                 |")
				//!PRINT "|   content is sufficient, no fertilization required                    |"
				//ELSE
			} else {
				//PRINT USING "|   ### kg N/ha   Duengung spaetestens bis ########                      |":DUNGbed*1.18,NAPPDAT$
				fmt.Printf("|   %03d kg N/ha   Duengung spaetestens bis %s                      |\n", int(g.DUNGBED*1.18), g.NAPPDAT)
				//!PRINT USING "|   ### kg N/ha   fertilizer latest by ##########                       |":DUNGbed*1.18,NAPPDAT$
				//END IF
			}
			// PRINT "*************************************************************************";
			fmt.Println("*************************************************************************")
			// LET ENDSTADIUM$ = "manueller Abbruch"
			g.ENDSTADIUM = invalidState
			// !LET ENDSTADIUM$ = "user break"
			// LET DUNG1 = DUNGBED*1.18
			//DUNG1 = g.DUNGBED * 1.18
			// LET DUNG2 = 0
			DUNG2 = 0
			// LET NAPPDAT2$ = "--------"
			NAPPDAT2 = "--------"
			// IF DUNGBED*1.18 < 2 THEN
			if g.DUNGBED*1.18 < 2 {
				//LET NAPPDAT$ = "--------"
				g.NAPPDAT = "--------"
				//END IF
			}
			//ELSE
		} else {
			//IF DUNGBED*1.18 < 65 then
			if g.DUNGBED*1.18 < 65 {
				//PRINT "*************************************************************************"
				fmt.Println("*************************************************************************")
				//PRINT USING "|   Prognose des Stickstoffbedarfs bis zum Stadium <###################  |":ENDSTADIUM$
				fmt.Printf("|   Prognose des Stickstoffbedarfs bis zum Stadium <%19s  |\n", g.ENDSTADIUM)
				//PRINT USING "|   Voraussichtliches Eintrittsdatum  ########                           |":ENDDAT$
				fmt.Printf("|   Voraussichtliches Eintrittsdatum  %s                           |\n", ENDDAT)
				//!PRINT USING "|   Prediction of nitrogen demand until stage      <################### |":ENDSTADIUM$
				//!PRINT USING "|   predicted date of stage         ##########                          |":ENDDAT$
				//IF DUNGBED*1.18 < 10 THEN
				if g.DUNGBED*1.18 < 10 {
					//PRINT "|   Vorrat ausreichend, keine Duengung notwendig.                        |"
					fmt.Println("|   Vorrat ausreichend, keine Duengung notwendig.                        |")
					//!  PRINT "|   content sufficient, no fertilizer required.                         |"
					// LET DUNG1 = 0
					DUNG1 = 0
					// LET DUNG2 = 0
					DUNG2 = 0
					// LET NAPPDAT$  = "--------"
					g.NAPPDAT = "--------"
					// LET NAPPDAT2$ = "--------"
					NAPPDAT2 = "--------"
					// ELSE IF DUNGBED*1.18 < 20 THEN
				} else if g.DUNGBED*1.18 < 20 {
					//PRINT USING "|   Empfohlene Duengung   20 kg N/ha   bis spaetestens ########          |":NAPPDAT$
					fmt.Printf("|   Empfohlene Duengung   20 kg N/ha   bis spaetestens %s          |\n", g.NAPPDAT)
					// !PRINT USING "|   recommended dosis    20 kg N/ha   by latest   ##########            |":NAPPDAT$
					// LET DUNG1 = 20
					DUNG1 = 20
					// LET DUNG2 = 0
					DUNG2 = 0
					// LET NAPPDAT2$ = "--------"
					NAPPDAT2 = "--------"
					//ELSE
				} else {
					//PRINT USING "|   Empfohlene D�uengung  ### kg N/ha   bis spaetestens ########          |":DUNGBED*1.18,NAPPDAT$
					fmt.Printf("|   Empfohlene Duengung  %03d kg N/ha   bis spaetestens %s          |\n", int(g.DUNGBED*1.18), g.NAPPDAT)
					// !PRINT USING "|   recommended dosis   ### kg N/ha   by latest    ##########           |":DUNGBED*1.18,NAPPDAT$
					// LET DUNG1 = DUNGBED*1.18
					// DUNG1 = g.DUNGBED * 1.18 //never use
					// LET DUNG2 = 0
					DUNG2 = 0
					// LET NAPPDAT2$ = "--------"
					NAPPDAT2 = "--------"
					//END IF
				}
				//PRINT "*************************************************************************";
				fmt.Println("*************************************************************************")
				//ELSE
			} else {
				//LET DUNG1 = (dungbed-dungbed/3)*1.18
				DUNG1 = (g.DUNGBED - (g.DUNGBED / 3)) * 1.18
				//IF DUNG1 > 65 THEN
				if DUNG1 > 65 {
					//LET DUNG1 = 65
					DUNG1 = 65
					//LET DUNG2 = DUNGBED*1.18-65
					DUNG2 = g.DUNGBED*1.18 - 65
					//ELSE
				} else {
					//LET dung2 = (dungbed-dungbed*2/3)*1.18
					DUNG2 = (g.DUNGBED - g.DUNGBED*2/3) * 1.18
					//END IF
				}
				//IF Endstadium$ = "Schossen" AND (W(1)+W(2)+W(3))/3 <  .3 then
				if g.ENDSTADIUM == schossen && (g.W[0]+g.W[1]+g.W[2])/3 < .3 {
					//Let du = dung1
					DU := DUNG1
					//LET DUNG1 = DUNG2
					DUNG1 = DUNG2
					//LET DUNG2 = DU
					DUNG2 = DU
					//END IF
				}
				// LET dunda2 = ROUND((ende-napp)/2+napp)
				DUNDA2 := int(math.Round(float64(g.ENDE-NAPP)/2 + float64(NAPP)))
				// If BLUET > 0 and Dunda2 > BLUET then let DUNDA2 = BLUET
				if g.BLUET > 0 && DUNDA2 > g.BLUET {
					DUNDA2 = g.BLUET
				}
				// LET LAST = ANJAHR*365+INT(ANJAHR/4)+165
				LAST := (g.ANJAHR)*365 + int(float64(g.ANJAHR)/4) + 165
				// IF DUNDA2 > LAST THEN let DUNDA2 = LAST
				if DUNDA2 > LAST {
					DUNDA2 = LAST
				}
				// CALL KALENDER(dunda2,nappdat2$)
				NAPPDAT2 = g.Kalender(DUNDA2)
				//       PRINT "*************************************************************************"
				fmt.Println("*************************************************************************")
				// PRINT USING "|   Prognose des Stickstoffbedarfs bis zum Stadium <###################  |":ENDSTADIUM$
				fmt.Printf("|   Prognose des Stickstoffbedarfs bis zum Stadium <%19s  |\n", g.ENDSTADIUM)
				// PRINT USING "|   Voraussichtliches Eintrittsdatum  ########   Duengung in 2 Gaben:    |":ENDDAT$
				fmt.Printf("|   Voraussichtliches Eintrittsdatum  %s   Duengung in 2 Gaben:    |\n", ENDDAT)
				// PRINT USING "|   ### kg N/ha   bis spaetestens ########,  ### kg N/ha   bis  ######## |":DUNG1,NAPPDAT$,dung2,nappdat2$
				fmt.Printf("|   %03d kg N/ha   bis spaetestens %s  %03d kg N/ha   bis  %s  |\n", int(DUNG1), g.NAPPDAT, int(DUNG2), NAPPDAT2)
				// !PRINT USING "|   Prediction of fertilizer demand until stage    <################### |":ENDSTADIUM$
				// !PRINT USING "|   predicted date of stage        ########## application in 2 doses:   |":ENDDAT$
				// !PRINT USING "|   ### kg N/ha   latest by    ##########,  ### kg N/ha   by ########## |":DUNG1,NAPPDAT$,dung2,nappdat2$
				// PRINT "*************************************************************************";
				fmt.Println("*************************************************************************")
				//END IF
			}
			//END IF
		}
		//!BOX KEEP 0,640,18,98 in RECOM$
		//!PRINT #8,USING "<####### ###.# ###.# ###.# ###.# ###.# ###.#":PROGDAT$,SUMA1,SUMA2,SUMA3,SUMA1+SUMA2+SUMA3,PLANA,OUTA
		// ELSE IF PROGNOS = ENDE THEN
	} else if g.PROGNOS == g.ENDE {
		// PRINT " aktuelle Nmin-Verteilung zum ";ENDDAT$;":"
		fmt.Println(" aktuelle Nmin-Verteilung zum " + ENDDAT + ":")
		// !PRINT " simulated Nmin distribution by ";ENDDAT$;":"
		// PRINT USING "                                             0 _ 30 cm: ### kg N/ha":SUMN1
		fmt.Printf("                                             0 _ 30 cm: %03d kg N/ha", int(SUMN1))
		// PRINT USING "                                            30 _ 60 cm: ### kg N/ha":SUMN2
		fmt.Printf("                                            30 _ 60 cm: %03d kg N/ha", int(SUMN2))
		// PRINT USING "                                            60 _ 90 cm: ### kg N/ha":SUMN3
		fmt.Printf("                                            60 _ 90 cm: %03d kg N/ha", int(SUMN3))
		// PRINT "                                            -----------------------"
		fmt.Println("                                            -----------------------")
		// PRINT USING "                                             0 _ 90 cm: ### kg N/ha":SUMN1+SUMN2+SUMN3
		fmt.Printf("                                             0 _ 90 cm: %03d kg N/ha", int(SUMN1+SUMN2+SUMN3))
		// PRINT
		// PRINT USING "Es wurden bereits ### kg N/ha durch die Pflanzen aufgenommen":PESUM
		fmt.Printf("Es wurden bereits %03d kg N/ha durch die Pflanzen aufgenommen", int(g.PESUM))
		// !PRINT USING "there are already ### kg N/ha taken up by crops":PESUM
		// LET DUNG1 = 0
		DUNG1 = 0
		// LET DUNG2 = 0
		DUNG2 = 0
		// LET NAPPDAT$ = "--------"
		g.NAPPDAT = "--------"
		// LET NAPPDAT2$ = "--------"
		NAPPDAT2 = "--------"
		// !PRINT #8,USING "<####### ###.# ###.# ###.# ###.# ###.# ###.#":ENDDAT$,SUMN1,SUMN2,SUMN3,SUMN1+SUMN2+SUMN3,PESUM,OUTSUM
		// ELSE
	} else {
		// PRINT "Programmlauf unterbrochen am ";ENDDAT$
		fmt.Println("Programmlauf unterbrochen am " + ENDDAT)
		// PRINT " aktuelle Nmin-Verteilung zum ";ENDDAT$;":"
		fmt.Println(" aktuelle Nmin-Verteilung zum " + ENDDAT + ":")
		// !PRINT "Program interrupted by ";ENDDAT$
		// !PRINT " simulated Nmin distribution by ";ENDDAT$;":"
		// PRINT USING "                                             0 _ 30 cm: ### kg N/ha":SUMN1
		fmt.Printf("                                             0 _ 30 cm: %03d kg N/ha\n", int(SUMN1))
		// PRINT USING "                                            30 _ 60 cm: ### kg N/ha":SUMN2
		fmt.Printf("                                            30 _ 60 cm: %03d kg N/ha\n", int(SUMN2))
		// PRINT USING "                                            60 _ 90 cm: ### kg N/ha":SUMN3
		fmt.Printf("                                            60 _ 90 cm: %03d kg N/ha\n", int(SUMN3))
		// PRINT "                                            -----------------------"
		fmt.Println("                                            -----------------------")
		// PRINT USING "                                             0 _ 90 cm: ### kg N/ha":SUMN1+SUMN2+SUMN3
		fmt.Printf("                                             0 _ 90 cm: %03d kg N/ha\n", int(SUMN1+SUMN2+SUMN3))
		// PRINT
		fmt.Println(" ")
		// PRINT USING "Es wurden bereits ### kg N/ha durch die Pflanzen aufgenommen":PESUM
		fmt.Printf("Es wurden bereits %03d kg N/ha durch die Pflanzen aufgenommen\n", int(g.PESUM))
		// !PRINT USING "there are already ### kg N/ha taken up by crops":PESUM
		// LET DUNG1 = 0
		DUNG1 = 0
		// LET DUNG2 = 0
		DUNG2 = 0
		// LET NAPPDAT$ = "--------"
		g.NAPPDAT = "--------"
		// LET NAPPDAT2$ = "--------"
		NAPPDAT2 = "--------"
		// LET ENDSTADIUM$ = "manueller Abbruch"
		g.ENDSTADIUM = invalidState
		// !LET ENDSTADIUM$ = "user break"
		// !PRINT #8,USING "<####### ###.# ###.# ###.# ###.# ###.# ###.#":ENDDAT$,SUMN1,SUMN2,SUMN3,SUMN1+SUMN2+SUMN3,PESUM,OUTSUM
		// END IF
	}
	//!PRINT #8,USING "<####### <################### ### <####### ### <#######":ENDDAT$,ENDSTADIUM$,DUNG1,NAPPDAT$,DUNG2,NAPPDAT2$
	//! ************** Tabellarischer Ausdruck der Düngungsempfehlung **************
	// LET fert$ = Path$ & "RESULT\D_" & STR$(slnr) & ".txt"
	fert := hPath.fert
	// OPEN #6:NAME fert$,ACCESS OUTIN,CREATE NEWOLD,ORGANIZATION TEXT
	// ERASE #6
	// CLOSE #6
	// OPEN #6:NAME fert$,ACCESS OUTIN,CREATE NEWOLD,ORGANIZATION TEXT
	// SET #6:MARGIN 90
	fertFile := OpenResultFile(fert, false)
	defer fertFile.Close()
	// PRINT #6:"Fläche Nr ";Slnr
	str := fmt.Sprintln("Fläche Nr " + strconv.Itoa(g.SLNR))
	if _, err := fertFile.Write(str); err != nil {
		log.Fatal(err)
	}
	// PRINT #6:"Datum         aktuelle Nmin - Werte  N-Düngung   düngen    Prognose    bis"
	str = fmt.Sprintln("Datum         aktuelle Nmin - Werte  N-Düngung   düngen    Prognose    bis")
	if _, err := fertFile.Write(str); err != nil {
		log.Fatal(err)
	}
	// PRINT #6:"            0-30 30-60 60-90  0-90cm  kg N/ha     bis         bis     Stadium "
	str = fmt.Sprintln("            0-30 30-60 60-90  0-90cm  kg N/ha     bis         bis     Stadium ")
	if _, err := fertFile.Write(str); err != nil {
		log.Fatal(err)
	}
	//LET FORM$ = "##########  ###   ###   ###    ###      ###    ########## ##########  <##################### "
	FORM := "%10s  %03d   %03d   %03d    %03d      %03d    %10s %10s  < %s "
	// PRINT #6,USING FORM$:PROGnoDAT$,SUMA1,SUMA2,SUMA3,SUMA1+SUMA2+SUMA3,ENDDAT$,DUNGBED*1.18,NAPPDAT$,Endstadium$
	str = fmt.Sprintf(FORM, progDatStr, int(SUMA1), int(SUMA2), int(SUMA3), int(SUMA1+SUMA2+SUMA3), int(g.DUNGBED*1.18), ENDDAT, g.NAPPDAT, g.ENDSTADIUM)
	if _, err := fertFile.Write(str); err != nil {
		log.Fatal(err)
	}
	// !PRINT #6:"ÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄÄ"
	// !print #6:" "
	// Let aktwas = 0
	// aktwas := 0.
	// // For I = 1 to n
	// for i := 0; i < g.N; i++ {
	// 	// let aktwas = aktwas + WG(1,i)
	// 	aktwas = aktwas + g.WG[1][i]
	// 	// next I
	// }

	// !print #8,using "Niederschlag ###.# ETakt ###.# Wgeh.nd. ####.# Sicker ####.## Bilanz ####.#":Regensum*10,evasum*10,(aktwas-ANFwas)*100,SCHNEE,(regensum-evasum-(aktwas-anfwas)*10)*10-SCHNEE
	// ! LET fert$ = Path$ & "RESULT\FERTreco.txt"
	// ! OPEN #6:NAME fert$,ACCESS OUTIN,CREATE NEWOLD,ORGANIZATION TEXT
	// ! set #6:pointer end
	// ! let appdat$ = NAPPDAT$(1:2) & NAPPDAT$(4:5) & NAPPDAT$(9:10)
	// ! PRINT #6, using "%%%%%     %%% KAS ######":slnr,dung1,APPDAT$
	// ! IF DUNG2 > 0 then
	// !  let appdat2$ = NAPPDAT2$(1:2) & NAPPDAT2$(4:5) & NAPPDAT2$(9:10)
	// !  PRINT #6, using "%%%%%     %%% KAS ######":slnr,dung2,APPDAT2$
	// ! end if
	//         close #6
	// DO
	//            CALL TC_Event (timer, event$, window, x1, x2)

	//            SELECT CASE event$
	//            CASE "KEYPRESS"
	//                 IF x1 = 27 THEN
	//                    EXIT DO        ! Exit the event loop when Escape pressed
	//                 END IF

	//            CASE "HIDE"
	//                 !IF window = 0 then
	//                 EXIT DO           ! Alternate exit
	//                 !END IF

	//                 ! Window-related events
	//            CASE "MENU"            ! x1 = menu number, x2 = menu item

	//                 ! Control-related events.  x2 = control id
	//            CASE "CONTROL SINGLE"
	//            CASE "CONTROL DOUBLE"
	//            CASE "CONTROL SELECT"

	//            CASE "CONTROL DESELECTED"

	//                 IF x2 = pbid_7 then
	//                    EXIT DO
	//                 END IF

	//            CASE ELSE

	//            END SELECT
	//         LOOP

	//         CALL TC_free(win_23)
	//     END SUB
}

// END MODULE
