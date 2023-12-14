#' @title Getting a hermes2go wrapper options list with initialized fields
#'
#' @description This function returns a default options list
#'
#' @param hermes2go_path Path of the binary executable file
#'
#' @param hermes2go_projects Path where to find the projects folder
#'
#' @param concurrency Number of concurrent simulations to run
#' (default: 1)
#'
#' @param time_display Logical value used to display (TRUE) or not (FALSE)
#' simulations duration
#'
#' @param warning_display Logical value used to display (TRUE) or not (FALSE)
#' (default: TRUE)
#'
#' @param weather_path Path where to find the weather files
#' (default: NULL)
#'
#' @param situation_parameters List containing the mapping between the situations names as data.frame
#'
#' @param out_path Path where to store the output files
#' (default: NULL)
#'
#' @param out_variable_names Vector of variable names to be returned.
#' If not provided, all variables will be returned.
#' (default: NULL)
#'
#' @return A list containing hermes2go wrapper options
#'
#' @examples
#' @export
#'
hermes2go_wrapper_options <- function(hermes2go_path,
                                      hermes2go_projects, ...) {
  # Template list
  options <- list()
  options$hermes2go_path <- character(0) # path
  options$hermes2go_projects <- character(0) # path
  options$concurrency <- 1 # integer
  options$time_display <- FALSE # boolean
  options$warning_display <- TRUE # boolean
  options$weather_path <- character(0) # path
  options$situation_parameters <- NULL # parameters for each situation (data.frame)
  options$out_path <- character(0) # path
  options$out_variable_names <- character(0) # vector of variable names


  # For getting the template
  # running hermes2go_wrapper_options
  if (!nargs()) {
    return(options)
  }

  # For fixing mandatory fields values
  options$hermes2go_path <- hermes2go_path
  options$hermes2go_projects <- hermes2go_projects

  # Fixing optional fields,
  # if corresponding to exact field names
  # in options list
  list_names <- names(options)
  add_args <- list(...)

  for (n in names(add_args)) {
    if (n %in% list_names) {
      options[[n]] <- add_args[[n]]
    }
  }

  return(options)
}

situation_params_from_excel <- function(excel_file) {
  # check if the excel file exists
  if (!file.exists(excel_file)) {
    stop("The excel file doesn't exist !")
  }
  # check if library is installed
  library(readxl)
  # read excel file
  situation_parameters <- readxl::read_excel(excel_file)

  # clean the situation parameters from execl overhead
  situation_parameters <- situation_parameters[!is.na(situation_parameters$SituationName), ]
  # remove class
  situation_parameters <- as.data.frame(situation_parameters)

  # check if the excel file contains the right columns
  if (!all(c("SituationName", "project", "plotNr") %in% colnames(situation_parameters))) {
    stop("The excel file should contain the following column: SituationName, project, plotNr")
  }

  # other columns can be overwrites of hermes Config struct parameters (see hermes/config.go)
  # check for hermes Config struct parameters
  valid_config_overwrites <- c(
    "Dateformat",
    "DivideCentury",
    "GroundWaterFrom",
    "ResultFileFormat",
    "ResultFileExt",
    "OutputIntervall",
    "ManagementEvents",
    "InitSelection",
    "SoilFile",
    "SoilFileExtension",
    "CropFileFormat",
    "MeasurementFileFormat",
    "PolygonGridFileName",
    "WeatherFile",
    "WeatherFileFormat",
    "WeatherNoneValue",
    "WeatherNumHeader",
    "CorrectionPrecipitation",
    "AnnualAverageTemperature",
    "ETpot",
    "CO2method",
    "CO2concentration",
    "CO2StomataInfluence",
    "NDeposition",
    "StartYear",
    "EndDate",
    "AnnualOutputDate",
    "VirtualDateFertilizerPrediction",
    "Latitude",
    "Altitude",
    "CoastDistance",
    "PTF",
    "LeachingDepth",
    "OrganicMatterMineralProportion",
    "KcFactorBareSoil",
    "PotMineralisation",
    "GroundWaterPhase",
    "Fertilization",
    "AutoSowingHarvest",
    "AutoFertilization",
    "AutoIrrigation",
    "AutoHarvest"
  )

  # other valid columns
  valid_default_columns <- c("SituationName", "project", "plotNr", "poligonID", "parameter", "gwId", "soilId", "fcode", "fileExtension")
  # project - folder name of the project
  # soilId - soil id
  # fcode - weather file code
  # fileExtension - file extention to automan, crop etc
  # plotNr - polyg in polygon file
  # poligonID - id extention to the plotNr in output file name
  # parameter - parameter folder overwrite (path)
  # gwId - groundwater id (if running with measured groundwater file)

  # remove columns that are not valid
  situation_parameters <- situation_parameters[, colnames(situation_parameters) %in% c(valid_config_overwrites, valid_default_columns)]

  # return the situation parameters
  return(situation_parameters)
}

#' @title funtion to check model options
#' @description This function checks the model options for Hermes2Go
#' @param model_options List containing any information needed by the model.
check_model_options <- function(model_options) {
  valid <- TRUE

  # check if path options are provided
  if (is.null(model_options$hermes2go_path) || is.null(model_options$hermes2go_projects)) {
    # print an error message
    print("hermes2go_path and hermes2go_projects should be elements of the model_model_options")
    valid <- FALSE
  } else {
    # check if the path exists
    if (!file.exists(model_options$hermes2go_path)) {
      print(paste("hermes2go_path doesn't exist !", model_options$hermes2go_path))
      valid <- FALSE
    } else {
      # check hermes version executable
      cmd <- paste(model_options$hermes2go_path, "--version")
      val <- system(cmd, wait = TRUE, intern = TRUE)
      if (!is.null(attr(val, "status"))) {
        print(paste(model_options$hermes2go_path, "is not executable or is not a hermes2go executable !"))
        valid <- FALSE
      }
    }
    # check if the project path exists
    if (!file.exists(model_options$hermes2go_projects)) {
      print(paste("hermes2go_projects doesn't exist !", model_options$hermes2go_projects))
      valid <- FALSE
    }
  }


  return(valid)
}
