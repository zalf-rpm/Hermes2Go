# check if package is installed
if (!requireNamespace("testthat", quietly = TRUE)) {
  install.packages("testthat")
}
library(testthat)

test_file("src/r_wrapper/testing/hermes_generate_batch_test.R")
test_file("src/r_wrapper/testing/hermes_wrapper_options_test.R")
test_file("src/r_wrapper/testing/hermes_wrapper_test.R")
