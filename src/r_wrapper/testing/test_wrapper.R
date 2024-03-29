#' @title Test model wrappers
#'
#' @description This function perform some tests of CroptimizR model wrappers.
#' See @details for more information.
#'
#' @param model_function Crop Model wrapper function to use.
#'
#' @param model_options List of options for the Crop Model wrapper (see help of
#' the Crop Model wrapper function used).
#'
#' @param param_values a named vector that contains values and names for AT LEAST TWO model parameters THAT ARE EXPECTED TO PLAY ON ITS RESULTS.
#'
#' @param sit_names Vector of situations names for which results must be tested.
#'
#' @param var_names (optional) Vector of variables names for which results must be tested. If not provided
#'
#' @details This function runs the wrapper consecutively with different subsets of param_values.
#' It then checks:
#'    - that its results are different when different subsets of param_values are used,
#'    - that its results are identical when same subsets of param_values are used.
#'
#' @return A list containing:
#'     - param_values_1: first subset of param_values
#'     - param_values_2: second subset of param_values
#'     - sim_1: results obtained with param_values_1
#'     - sim_2: results obtained with param_values_2
#'     - sim_3: results obtained for second run with param_values_1
#'
test_wrapper <- function(model_function, model_options, param_values, sit_names, var_names = NULL) {
  if (length(param_values) <= 1) {
    stop("param_values argument must include at least TWO parameters.")
  }

  param_values_1 <- param_values[1]
  param_values_2 <- param_values
  sim_1 <- model_function(
    param_values = param_values_1, model_options = model_options,
    sit_names = sit_names, var_names = var_names
  )
  sim_2 <- model_function(
    param_values = param_values_2, model_options = model_options,
    sit_names = sit_names, var_names = var_names
  )
  sim_3 <- model_function(
    param_values = param_values_1, model_options = model_options,
    sit_names = sit_names, var_names = var_names
  )

  results <- list(
    param_values_1 = param_values_1,
    param_values_2 = param_values_2,
    sim_1 = sim_1,
    sim_2 = sim_2,
    sim_3 = sim_3
  )

  cat(crayon::green("Test the wrapper gives identical results when running with same inputs ...\n"))
  if (identical(sim_1, sim_3)) {
    cat(crayon::green("... OK\n"))
  } else {
    cat(crayon::red("... test failed\n"))
    cat(crayon::red("The wrapper gave different results although executed two times \n 
    with the same inputs param_values_1 (see results in sim_1 and sim_3).\n"))
    cat(crayon::red("This may be due to model input files that do not come back to their original state at the end of the wrapper.\n"))
  }
  cat("\n")
  cat(crayon::green("Test the wrapper gives different results when running with different inputs ...\n"))
  if (!identical(sim_1, sim_2)) {
    cat(crayon::green("... OK\n"))
  } else {
    cat(crayon::red("... test failed\n"))
    cat(crayon::red("The wrapper gave same results although executed with different param_values,
     arguments param_values_1 and param_values_2 (see results in sim_1 and sim_2).\n"))
    cat(crayon::red("Either param_values is not correctly handled in the model wrapper 
    (in part. make sure that model inputs are modified according to param_values) 
    or the selected parameters do not play on the model outputs (test with other parameters in param_values).\n"))
  }

  return(invisible(results))
}
