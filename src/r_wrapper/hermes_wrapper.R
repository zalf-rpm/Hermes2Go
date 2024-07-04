#' @title Running Hermes2Go from txt input files stored in one directory
#' per `situation`, simulated results are returned in a list
#'
#' @description This function uses Hermes2Go through a system call to run.
#' It requires a valid Hermes2Go executable file and
#' a directory containing the project(s)
#'
#' @param param_values (optional) a named vector that contains the value(s) and name(s)
#' of the parameters to force for each situation to simulate. If not provided (or if is
#' NULL), the simulations will be performed using default values of the parameters
#' (e.g. as read in the model input files).
#'
#' @param model_options List containing any information needed by the model.
#' Use hermes2go_wrapper_options to get a template list with initialized fields.
#' hermes2go_path - the path of Hermes2Go executable file
#' hermes2go_projects - path of the directory containing the Hermes2Go input data
#' concurrency (optional) - number of parallel processes to run
#' time_display (optional) - if TRUE, the function will display the time taken to run the model
#' warning_display (optional) - if TRUE, the function will display the warnings
#' weather_path (optional) - the path of the directory containing the weather
#' files (if not provided, the weather folder needs to be provided in project config file)
#'
#' @param sit_names Vector of situations names for which results must be returned.
#' a situation should match a directory name in the hermes2go_projects directory
#'
#' @param var_names (optional) Vector of variable names to be returned.
#' If not provided, all variables will be returned.
#'
#' @return A list containing simulated values (`sim_list`: a vector of list (one
#' element per values of parameters) containing usms outputs data.frames) and an
#' error code (`error`) indicating if at least one simulation ended with an
#' error.
#'
#' @examples
#' @export
hermes2go_wrapper <- function(param_values,
                              model_options,
                              sit_names = NULL,
                              var_names = NULL,
                              ...) {

  # check if all the required options are provided
  if (!check_model_options(model_options)) {
    stop("Invalid model options")
  }

  hermes2go_path <- model_options$hermes2go_path # path
  hermes2go_projects <- model_options$hermes2go_projects # path
  concurrency <- model_options$concurrency # integer
  warning_display <- model_options$warning_display # boolean
  weather_path <- model_options$weather_path # path
  result_folder <- model_options$out_path # path
  situation_parameters <- model_options$situation_parameters # path
  use_temp_dir <- model_options$use_temp_dir # boolean
  output_function <- model_options$output_function # function
  crop_file_name <- model_options$crop_file_name # string

  # check if param_values is an array, get the number of rows
  num_rows <- 1
  if (base::is.array(param_values)) {
    num_rows <- dim(param_values)[1]
  }
  # results
  res <- list()
  res$error <- FALSE
  res$sim_list <- list()

  # TODO: need to implement crop parameter replacement in hermes2go

  # track execution time
  start_time <- Sys.time()

  # if no result folder provided, use temp dir
  if (is.null(result_folder)) {
    result_folder <- tempdir()
  }
  # create sub dir in temp dir
  result_folder <- file.path(result_folder, "hermes2go_results")

  # delete temp dir if it already exists, to make sure we don't have old results
  if (dir.exists(result_folder)) {
    unlink(result_folder, recursive = TRUE)
  }
  dir.create(result_folder)

  batch_file <- generate_batch_file(param_values, sit_names, situation_parameters, crop_file_name, weather_path, result_folder, use_temp_dir)

  # Run Herme2Go ------------------------------------------------------------------
  cmd <- paste(hermes2go_path,
    sep = " ", collapse = "",
    "-module batch", "-batch", batch_file,
    "-concurrent", concurrency,
    "-workingdir", hermes2go_projects
  )
  # Display the command to run ---------------------------------------------------
  if (warning_display) {
    print(paste("Running Hermes2Go with command:", cmd))
  }
  # run Hermes2Go
  run_file_stdout <- system(cmd, wait = TRUE, intern = TRUE)

  # Getting the execution status
  res$error <- !is.null(attr(run_file_stdout, "status"))

  # TODO: log the output of Hermes2Go in a file
  if (res$error) {
    print(run_file_stdout)
  }

  for (ip in 1:num_rows) {
    # Store results ---------------------------------------------------------------
    results_tmp <- output_function(
      result_folder,
      sit_names,
      var_names,
      param_values
    )
    res$sim_list <- results_tmp
  }

  # Display simulation duration -------------------------------------------------
  if (model_options$time_display) {
    duration <- Sys.time() - start_time
    print(duration)
  }

  if (length(res$sim_list) > 0) {
    # Add the attribute cropr_simulation for using CroPlotR package
    attr(res$sim_list, "class") <- "cropr_simulation"
  }
  if (dir.exists(result_folder)) {
    res$result_folder <- result_folder
  }
  return(res)
}
