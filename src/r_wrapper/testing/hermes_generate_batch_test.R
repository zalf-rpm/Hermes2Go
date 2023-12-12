
test_that("situation parameters are converted into lines", {
  # Load the function under test
  source(file.path("..", "hermes_generate_batch.R"))

  # Define test inputs
  sit_names <- c("sit1", "sit2")

  # generate situation parameters with SituationName Parameter1 Parameter2 as columns and 4 rows
  situation_parameters <- data.frame(
    "SituationName" = c("sit1", "sit2", "sit3", "sit4"),
    "Parameter1" = c(10, 20, 30, 40),
    "Parameter2" = c(100, 200, 300, 400)
  )

  # Generate the expected lines
  expected_lines <- list(
    "sit1" = "Parameter1=10 Parameter2=100",
    "sit2" = "Parameter1=20 Parameter2=200"
  )
  print(expected_lines)

  # Call the function under test
  actual_lines <- situation_parameters_to_line(sit_names, situation_parameters)

  print(actual_lines)

  # Compare the actual and expected lines
  expect_equal(actual_lines, expected_lines)
})


test_that("generate_batch_file generates the correct batch file", {

  # Load the function under test
  source(file.path("..", "hermes_generate_batch.R"))

  # Define test inputs
  param_values <- c("CropFile" = "PARAM_varityX.WW", "TSum1" = 1, "TSum2" = 2)
  sit_names <- c("sit1", "sit2")
  situation_parameters <- data.frame(
    "Situation Name" = c("sit1", "sit2"),
    "Parameter1" = c(10, 20),
    "Parameter2" = c(100, 200)
  )
  weather_path <- "/path/to/weather"
  result_folder <- "/path/to/results"

  # Generate the expected batch file
  expected_batch_file <- tempfile("expected_batch", fileext = ".txt")
  expected_lines <- c(
    "Situation Name=sit1 Parameter1=10 Parameter2=100 CropFile=PARAM_varityX.WW c_TSum_1=1 c_TSum_2=2 WeatherRootFolder=/path/to/weather resultfolder=/path/to/results/sit1", # nolint: line_length_linter.
    "Situation Name=sit2 Parameter1=20 Parameter2=200 CropFile=PARAM_varityX.WW c_TSum_1=1 c_TSum_2=2 WeatherRootFolder=/path/to/weather resultfolder=/path/to/results/sit2" # nolint: line_length_linter.
  )
  writeLines(expected_lines, expected_batch_file)

  # Call the function under test
  actual_batch_file <- generate_batch_file(param_values, sit_names, situation_parameters, weather_path, result_folder)

  # Compare the actual and expected batch files
  expect_equal(readLines(actual_batch_file), readLines(expected_batch_file))
})
