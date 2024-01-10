test_that("hermes wrapper runs with calibration example", {
  if (!requireNamespace("readxl", quietly = TRUE)) {
    install.packages("readxl")
  }

  source(file.path("..", "hermes_wrapper_options.R"))
  source(file.path("..", "hermes_generate_batch.R"))
  source(file.path("..", "hermes_wrapper.R"))

  # Define test inputs
  excel_file <- file.path("data", "situation_parameter.xlsx")
  print(excel_file)
  situation_parameters <- situation_params_from_excel(excel_file)

  # path to exe
  hermes2go_path <- file.path("../..", "hermes2go", "hermes2go.exe")
  hermes2go_path <- normalizePath(hermes2go_path)
  # check if the exe file exists
  if (!file.exists(hermes2go_path)) {
    stop("The exe file doesn't exist !")
  }

  hermes2go_projects <- file.path("../../..", "examples")
  weather_path <- file.path("../../..", "examples", "weather")
  out_path <- file.path("../../..", "examples", "calibration_output")

  # check if the output folder exists
  if (!dir.exists(out_path)) {
    dir.create(out_path)
  }
  out_path <- normalizePath(out_path)

  var_names <- c("Crop", "Year", "Yield", "MaxLAI", "SowDOY", "sum_ET", "sum_irri", "AWC_30_sow", "AWC_30_harv")

  model_options <- hermes2go_wrapper_options(hermes2go_path,
                                             hermes2go_projects,
                                             concurrency = 2,
                                             time_display = TRUE,
                                             weather_path = weather_path,
                                             situation_parameters = situation_parameters,
                                             out_path = out_path)

  # situation vector
  sit_names <- c("sit2", "sit3", "sit4")
  param_values <- c(0.5, 0.6, 0.7)

  # Call the function under test
  result <- hermes2go_wrapper(param_values, model_options, sit_names, var_names)

})