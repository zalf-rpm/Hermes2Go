test_that("hermes wrapper runs with calibration example", {
  if (!requireNamespace("readxl", quietly = TRUE)) {
    install.packages("readxl")
  }

  source(file.path("..", "hermes_wrapper_options.R"))
  source(file.path("..", "hermes_generate_batch.R"))
  source(file.path("..", "hermes_wrapper.R"))

  # test inputs
  excel_file <- file.path("data", "situation_parameter.xlsx")
  situation_parameters <- situation_params_from_excel(excel_file)

  # path to exe
  hermes2go_path <- file.path("../..", "hermes2go", "hermes2go.exe")
  hermes2go_path <- normalizePath(hermes2go_path)
  # check if the exe file exists
  if (!file.exists(hermes2go_path)) {
    stop("The exe file doesn't exist !")
  }
  # path to projects
  hermes2go_projects <- file.path("../../..", "examples")
  hermes2go_projects <- normalizePath(hermes2go_projects)
  # path to weather data
  weather_path <- file.path("../../..", "examples", "weather")
  weather_path <- normalizePath(weather_path)
  # path to output
  out_path <- file.path("../../..", "examples", "calibration_output")

  # check if the output folder exists
  if (!dir.exists(out_path)) {
    dir.create(out_path)
  }
  out_path <- normalizePath(out_path)

  # output variable filter
  var_names <- c("Crop", "Year", "Yield", "MaxLAI", "SowDOY", "sum_ET", "sum_irri", "AWC_30_sow", "AWC_30_harv")

  # model options generation
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
  # calibration parameters
  param_values <- list(
    "CropFile" = "PARAM.SOY",
    "c_TSUM_1" = 73,
    "c_TSUM_2" = 55,
    "c_TSUM_3" = 240,
    "c_TSUM_4" = 330
  )

  # Call the function under test
  result <- hermes2go_wrapper(param_values, model_options, sit_names, var_names)
  print(result)
  expected_result <- list(
    "error" = FALSE,
    "sim_list" = list(
      "sit2" = list(
        "Csit210002" = data.frame(
          "Crop" = "SOY",
          "Year" = 2002,
          "Yield" = 3599,
          "MaxLAI" = 9.0,
          "SowDOY" = 135,
          "sum_ET" = 44,
          "sum_irri" = 0,
          "AWC_30_sow" = 115,
          "AWC_30_harv" = 120
        ),
        "Vsit210002" = data.frame()
      ),
      "sit3" = list(
        "Csit310003" = data.frame(
          "Crop" = "SM ",
          "Year" = 2003,
          "Yield" = 4701,
          "MaxLAI" = 1.1,
          "SowDOY" = 135,
          "sum_ET" = 56,
          "sum_irri" = 0,
          "AWC_30_sow" = 83,
          "AWC_30_harv" = 31
        ), 
        "Vsit310003" = data.frame()
      ),
      "sit4" = list(
        "Csit410004" = data.frame(
          "Crop" = "SOY",
          "Year" = 2003,
          "Yield" = 2650,
          "MaxLAI" = 6.3,
          "SowDOY" = 135,
          "sum_ET" = 49,
          "sum_irri" = 0,
          "AWC_30_sow" = 106,
          "AWC_30_harv" = 84
        ),
        "Vsit410004" = data.frame()
      )

    )
  )
  attr(expected_result$sim_list, "class") <- "cropr_simulation"
  expected_result$result_folder <- file.path(out_path, "hermes2go_results")

  expect_equal(result, expected_result)

})