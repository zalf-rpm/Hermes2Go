#' @title generate hermes batch file for running simulations
#'
#' @description This function generates a batch file which can be passed
#' to Hermes2Go in order to apply parameter value changes using parameters,
#' a named vector of values.
#'
#' @param param_values a named vector of parameters values
#'
#' @param sit_names Vector of situations names for which results must be returned.
#'
#' @param situation_parameters list containing the mapping between the situations names
#' and the project folders and settings
#'
#' @param crop_file_name Name of the crop file to use for calibrations
#' 
#' @param weather_path Path where to find the root folder for weather files
#'
#' @param result_folder Path where to store the output files
#'
#' @return path to batch file
#'
#' @export
#'
generate_batch_file <- function(param_values, sit_names, situation_parameters, crop_file_name, weather_path, result_folder, use_temp_dir = TRUE) {

  # if result_folder is null, use the current folder
  if (is.null(result_folder)) {
    result_folder <- ""
  }

  if (use_temp_dir) {
    batch_file <- tempfile("hermes_batch", fileext = ".txt")
  } else {
    batch_file <- file.path(result_folder, "hermes_batch.txt")
  }
  file_conn <- file(batch_file)

  if (!is.null(weather_path)) {
    weather <- paste("WeatherRootFolder", weather_path, sep = "=")
  }
  crop_param <- ""
  if (!is.null(crop_file_name)) {
    crop_param <- paste("CropFile", crop_file_name, sep = "=")
  }

  sim_lines <- situation_parameters_to_line(sit_names, situation_parameters)
  param_lines <- params_to_line(param_values)

  num_sims <- length(sim_lines)
  num_param_lines <- length(param_lines)
  # to run without parameters
  if (num_param_lines == 0) {
    num_param_lines <- 1
  }
  lines <- vector("character", num_param_lines * num_sims)

  # combine each simulation line with each parameter line
  ip <- 0
  for (id in seq_along(sim_lines)) {

    # get key of sim
    sim_name <- names(sim_lines)[id]
    # get value of sim
    sim <- sim_lines[[id]]

    sim_result_folder <- paste(result_folder, sim_name, sep = "/")
    sim_result_folder <- paste("resultfolder", sim_result_folder, sep = "=")

    for (param in param_lines) {
      ip <- ip + 1
      lines[ip] <- paste(sim, crop_param, param, weather, sim_result_folder, sep = " ")
    }
  }

  writeLines(lines, file_conn)
  close(file_conn)
  return(batch_file)
}

situation_parameters_to_line <- function(sit_names, situation_parameters) {
  # check if situation_parameters has a `SituationName` attribute
  if (!"SituationName" %in% names(situation_parameters)) {
    print(names(situation_parameters))
    stop("SituationName must be provided as a column in situation_parameters")
  }
  if (is.null(sit_names)) {
    # if no situation names provided, use all the situations
    filtered_situation_parameters <- situation_parameters
  } else {
    # filter by attribute `SituationName` the situation_parameters list to keep only the situations to simulate
    filtered_situation_parameters <- situation_parameters[situation_parameters$`SituationName` %in% sit_names, ]
  }
  # number of simulation rows
  sim_num <- dim(filtered_situation_parameters)[1]
  sim_lines <- list()

  # for each simulation row, generate a line of arguments

  for (is in 1:sim_num) {
    # get row at index is
    sim <- filtered_situation_parameters[is, , drop = FALSE]

    sim_line <- ""
    sim_name <- ""
    for (i in seq_along(sim)) {
      # must have a `Situation Name` attribute
      hermes_arg <- ""
      if (names(sim)[i] == "SituationName") {
        sim_name <- as.character(sim[i])
      } else {
        hermes_arg <- paste(names(sim)[i], as.character(sim[i]), sep = "=")
        if (sim_line == "") {
          sim_line <- hermes_arg
        } else {
          sim_line <- paste(sim_line, hermes_arg, sep = " ")
        }
      }
    }
    if (sim_name == "") {
      stop("Situation Name must be provided")
    }
    sim_lines[sim_name] <- sim_line
  }
  return(sim_lines)
}

# transform parameters values to a line of arguments
params_to_line <- function(param_values) {
  # check if param_values not null
  if (is.null(param_values)) {
    return(character(0))
  }
  parameter_names <- names(param_values)

  line <- ""
  for (ip in seq_along(param_values)) {
    param_name <- parameter_names[ip]
    param_value <- param_values[ip]
    # check if parameter name is not null
    if (is.null(param_name)) {
      stop("Parameter name must be provided")
    }
    # convert parameter name to Hermes2Go name
    translated_name <- predefinded_agmip_params(param_name)
    hermes_arg <- paste(translated_name, as.character(param_value), sep = "=")
    if (line == "") {
      line <- hermes_arg
    } else {
      line <- paste(line, hermes_arg, sep = " ")
    }
  }
  return(line)
}

# mapping between Apgmip and Hermes2Go parameters
predefinded_agmip_params <- function(param) {
  # list with Agmip name to Hermes2Go name mapping
  crop_overwrite_params <- list(
    "TSum1" = "c_TSum_1",
    "TSum2" = "c_TSum_2",
    "TSum3" = "c_TSum_3",
    "Vern2" = "c_Vern_2",
    "DLen2" = "c_DLen_2",
    "DLen3" = "c_DLen_3",
    "SLA1" = "c_SLA_1",
    "SLA2" = "c_SLA_2",
    "SLA3" = "c_SLA_3",
    "SLA4" = "c_SLA_4",
    "MaxEffectRootDepth" = "c_MaxEffectRootDepth",
    "AMAX" = "c_AMAX",
    "Kc4" = "c_Kc4",
    "part2org2_shoot" = "c_part2org_2_shoot",
    "part2org3_shoot" = "c_part2org_3_shoot",
    "part2org4_shoot" = "c_part2org_4_shoot",
    "ConcN_biom" = "cConcN_biom",
    "ConcN_root" = "cConcN_root",
    "N_root3" = "c_N_root_3",
    "N_root4" = "c_N_root_4",
    "percOrgan" = "c_percOrgan",
    "OrgMatterMinProp" = "OrganicMatterMineralProportion"
  )
  config_overwrite_params <- list(
    "OrgMatterMinProp" = "OrganicMatterMineralProportion"
  )

  if (param %in% names(crop_overwrite_params)) {
    return(crop_overwrite_params[[param]])
  } else if (param %in% names(config_overwrite_params)) {
    return(config_overwrite_params[[param]])
  } else {
    return(param)
  }
}
