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
  hermes2go_projects <- normalizePath(hermes2go_projects)
  weather_path <- file.path("../../..", "examples", "weather")
  weather_path <- normalizePath(weather_path)
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
                                             warning_display = TRUE,
                                             weather_path = weather_path,
                                             situation_parameters = situation_parameters,
                                             out_path = out_path,
                                             use_temp_dir = FALSE)

  # situation vector
  sit_names <- c("sit2", "sit3", "sit4")
  param_values <- list(
    "CropFile" = "PARAM_0.SOY",
    "TSum1" = 73,
    "TSum2" = 55,
    "TSum3" = 240,
    "TSum4" = 330
  )


  # Call the function under test
  result <- hermes2go_wrapper(param_values, model_options, sit_names, var_names)
  expected_result <- list(
    "error" = FALSE,
    "sim_list" = list(
      "sit2" = list(
        "Crop" = "SOY",
        "Year" = 2001,
        "Yield" = 3.5,
        "MaxLAI" = 3.5,
        "SowDOY" = 120,
        "sum_ET" = 0,
        "sum_irri" = 0,
        "AWC_30_sow" = 0,
        "AWC_30_harv" = 0
      ),
      "sit3" = list(
        "Crop" = "SOY",
        "Year" = 2002,
        "Yield" = 3.5,
        "MaxLAI" = 3.5,
        "SowDOY" = 120,
        "sum_ET" = 0,
        "sum_irri" = 0,
        "AWC_30_sow" = 0,
        "AWC_30_harv" = 0
      ),
      "sit4" = list(
        "Crop" = "SOY",
        "Year" = 2003,
        "Yield" = 3.5,
        "MaxLAI" = 3.5,
        "SowDOY" = 120,
        "sum_ET" = 0,
        "sum_irri" = 0,
        "AWC_30_sow" = 0,
        "AWC_30_harv" = 0
      )
    )
  )
  print(result)
  expect_equal(result, expected_result)

})