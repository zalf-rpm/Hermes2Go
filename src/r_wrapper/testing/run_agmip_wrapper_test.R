if (!requireNamespace("readxl", quietly = TRUE)) {
  install.packages("readxl")
}
if (!requireNamespace("here", quietly = TRUE)) {
  install.packages("here")
}

source(here::here("src/r_wrapper", "hermes_wrapper_options.R"))
source(here::here("src/r_wrapper", "hermes_generate_batch.R"))
source(here::here("src/r_wrapper", "hermes_wrapper.R"))
source(here::here("src/r_wrapper/testing", "test_wrapper.R"))

# test inputs
excel_file <- here::here("src/r_wrapper/testing/data", "situation_parameter.xlsx")
situation_parameters <- situation_params_from_excel(excel_file)

# path to exe
hermes2go_path <- here::here("src/hermes2go", "hermes2go.exe")
hermes2go_path <- normalizePath(hermes2go_path)
# check if the exe file exists
if (!file.exists(hermes2go_path)) {
  stop("The exe file doesn't exist !")
}
# path to projects
hermes2go_projects <- here::here("examples")
hermes2go_projects <- normalizePath(hermes2go_projects)
# path to weather data
weather_path <- here::here("examples/weather")
weather_path <- normalizePath(weather_path)
# path to output
out_path <- here::here("examples/calibration_output")

# check if the output folder exists
if (!dir.exists(out_path)) {
  dir.create(out_path)
}
out_path <- normalizePath(out_path)

# model options generation
model_options <- hermes2go_wrapper_options(hermes2go_path,
  hermes2go_projects,
  concurrency = 2,
  time_display = TRUE,
  warning_display = TRUE,
  weather_path = weather_path,
  situation_parameters = situation_parameters,
  out_path = out_path,
  use_temp_dir = FALSE
)

# situation vector
sit_names <- c("sit2", "sit3", "sit4")
# calibration parameters
param_values <- list(
  CO2concentration = 365,
  NDeposition = 10
)

test_wrapper(hermes2go_wrapper, model_options, param_values, sit_names, var_names = NULL)
