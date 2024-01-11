test_that("situation params from excel", {
  if (!requireNamespace("readxl", quietly = TRUE)) {
    install.packages("readxl")
  }
  # Load the function under test
  source(file.path("..", "hermes_wrapper_options.R"))

  # Define test inputs
  excel_file <- file.path("data", "situation_parameter.xlsx")
  print(excel_file)

  # Generate the expected lines
  expected_situation_parameters <- data.frame(
    "SituationName" = c("sit1", "sit2", "sit3", "sit4", "sit5", "sit6", "sit7", "sit8"),
    "project" = c("calibration", "calibration", "calibration", "calibration", "calibration", "calibration", "calibration", "calibration"),
    "plotNr" = c(10001, 10002, 10003, 10004, 10005, 10006, 10007, 10008),
    "poligonID" = c("sit1", "sit2", "sit3", "sit4", "sit5", "sit6", "sit7", "sit8"),
    "soilId" = c("002", "011", "002", "011", "002", "011", "075", "002"),
    "fcode" = c("109_120", "109_120", "109_120", "109_120", "109_120", "109_120", "109_120", "109_120"),
    "Altitude" = c(73, 73, 73, 73, 73, 73, 73, 73),
    "Latitude" = c("52.52", "52.52", "52.52", "52.52", "52.52", "52.52", "52.52", "52.52"), # is a string excel file
    "EndDate" = c(12312002, 12312002, 12312003, 12312003, 12312004, 12312004, 12312005, 12312005)
  )

  # SituationName	project	    plotNr	poligonID	soilId	fcode	  Altitude	Latitude	EndDate
  # sit1	        calibration	10001	  sit1	    002	    109_120	73	      52.52	    12312002
  # sit2	        calibration	10002	  sit2	    011	    109_120	73	      52.52	    12312002
  # sit3	        calibration	10001	  sit3	    002	    109_120	73	      52.52	    12312003
  # sit4	        calibration	10002	  sit4	    011	    109_120	73	      52.52	    12312003
  # sit5	        calibration	10001	  sit5	    002	    109_120	73	      52.52	    12312004
  # sit6	        calibration	10002	  sit6	    011	    109_120	73	      52.52	    12312004
  # sit7	        calibration	10001	  sit7	    075	    109_120	73	      52.52	    12312004
  # sit8	        calibration	10002	  sit8	    002	    109_120	73	      52.52	    12312005

  # Call the function under test
  actual_situation_parameters <- situation_params_from_excel(excel_file)

  # Compare the actual and expected lines
  expect_equal(actual_situation_parameters, expected_situation_parameters)
})