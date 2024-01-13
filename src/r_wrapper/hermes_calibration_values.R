list_crop_params <- function() {
  print("c_MAXAMAX     Amax (C-Assimilation at light saturation) at optimal temperature (kg CO2/ha leave/h)")
  print("c_MINTMP      Min temperatur for plant growth (°C)")
  print("c_WUMAXPF     crop specific root length (dm)")
  print("c_VELOC       root depth increase in mm/C°")
  print("c_YIFAK 	     percentage of harvested organ to yield (80%) ")
  print("c_INITCONCNBIOM initial N concentration in biomass (2%)")
  print("c_INITCONCNROOT initial N concentration in root (2%)")

  print("replace i by development stage number (1-6)")

  print("c_TSUM_<i> 	    temperature sum for stage i (°C)")
  print("c_BAS_<i> 	      base temperature for temperature sum (°C)")
  print("c_VSCHWELL_<i>   required vernalisation days for stage i (d)")
  print("c_DAYL_<i> 	    day length for stage i (+/-24 h)")
  print("c_DLBAS_<i> 	    base value for day length for stage i (+/-24h)")
  print("c_DRYSWELL_<i> 	dry stress threshold (Ta/Tp) for stage (0-1)")
  print("c_LUKRIT_<i> 	  critical air pore fraction for stage i (cm^3/cm^3)")
  print("c_LAIFKT_<i> 	  SLA specific leave area (area per mass) (m2/m2/kg TM) in stage i")
  print("c_WGMAX_<i> 	    N-content root end of stage i")
  print("c_KC_<i> 	      crop factor for evapotranspiration in stage i")

  print("replace j by partition number (1-4)")

  print("c_PRO_<i>_<j> 	    fraction of production for organ j in stage i (0-1)")
  print("c_DEAD_<i>_<j> 	  fraction of dead material for organ j in stage i (0-1)")
}
