package com.example.tempo.navigation

import android.app.Application
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.platform.LocalContext
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navigation
import com.example.tempo.TempoApp
import com.example.tempo.data.AuthRepository
import com.example.tempo.data.UserPreferencesRepository
import com.example.tempo.data.remote.RetrofitInstance
import com.example.tempo.ui.auth.AuthState
import com.example.tempo.ui.auth.AuthViewModel
import com.example.tempo.ui.auth.LoginScreen
import com.example.tempo.ui.auth.RegisterScreen

object AppRoutes {
    const val AUTH_GRAPH = "auth_graph"
    const val LOGIN = "login"
    const val REGISTER = "register"
    const val MAIN_APP = "main_app"
}

@Composable
fun AppNavigation() {
    val context = LocalContext.current
    val userPreferencesRepository = remember { UserPreferencesRepository(context) }
    val authToken by userPreferencesRepository.authToken.collectAsState(initial = null)
    val navController = rememberNavController()

    // This effect will react to changes in the authToken.
    // If the token is null, it means the user has logged out or the token expired.
    LaunchedEffect(authToken) {
        if (authToken == null) {
            navController.navigate(AppRoutes.AUTH_GRAPH) {
                // Clear the back stack to prevent going back to the main app
                popUpTo(AppRoutes.MAIN_APP) { inclusive = true }
            }
        }
    }

    val startDestination = if (authToken != null) AppRoutes.MAIN_APP else AppRoutes.AUTH_GRAPH

    NavHost(navController = navController, startDestination = startDestination) {
        authGraph(navController)

        composable(route = AppRoutes.MAIN_APP) {
            TempoApp()
        }
    }
}


fun NavGraphBuilder.authGraph(navController: NavHostController) {
    navigation(startDestination = AppRoutes.LOGIN, route = AppRoutes.AUTH_GRAPH) {
        composable(AppRoutes.LOGIN) {
            val viewModel: AuthViewModel = getAuthViewModel()
            val authState by viewModel.authState.collectAsState()

            // Navigate on success
            LaunchedEffect(authState) {
                if (authState is AuthState.Success) {
                    navController.navigate(AppRoutes.MAIN_APP) {
                        popUpTo(AppRoutes.AUTH_GRAPH) { inclusive = true }
                    }
                }
            }

            LoginScreen(
                authState = authState,
                onLoginClicked = viewModel::login,
                onNavigateToRegister = { navController.navigate(AppRoutes.REGISTER) },
                onDismissError = viewModel::resetState
            )
        }

        composable(AppRoutes.REGISTER) {
            val viewModel: AuthViewModel = getAuthViewModel()
            val authState by viewModel.authState.collectAsState()

            // Navigate on success
            LaunchedEffect(authState) {
                if (authState is AuthState.Success) {
                    navController.navigate(AppRoutes.MAIN_APP) {
                        popUpTo(AppRoutes.AUTH_GRAPH) { inclusive = true }
                    }
                }
            }

            RegisterScreen(
                authState = authState,
                onRegisterClicked = viewModel::register,
                onNavigateToLogin = { navController.popBackStack() },
                onDismissError = viewModel::resetState
            )
        }
    }
}


@Composable
private fun getAuthViewModel(): AuthViewModel {
    val context = LocalContext.current
    val factory = remember {
        object : ViewModelProvider.Factory {
            override fun <T : ViewModel> create(modelClass: Class<T>): T {
                val application = context.applicationContext as Application
                val userPreferencesRepository = UserPreferencesRepository(application)
                val authRepository = AuthRepository(RetrofitInstance.api)

                @Suppress("UNCHECKED_CAST")
                return AuthViewModel(authRepository, userPreferencesRepository) as T
            }
        }
    }
    return viewModel(factory = factory)
}