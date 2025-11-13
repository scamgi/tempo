package com.example.tempo.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.tempo.data.AuthRepository
import com.example.tempo.data.UserPreferencesRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class AuthViewModel(
    private val authRepository: AuthRepository,
    private val userPreferencesRepository: UserPreferencesRepository
) : ViewModel() {

    private val _authState = MutableStateFlow<AuthState>(AuthState.Idle)
    val authState: StateFlow<AuthState> = _authState.asStateFlow()

    fun login(email: String, password: String) {
        viewModelScope.launch {
            _authState.value = AuthState.Loading
            val result = authRepository.login(email, password)
            result.onSuccess { token ->
                userPreferencesRepository.saveAuthToken(token)
                _authState.value = AuthState.Success
            }.onFailure {
                _authState.value = AuthState.Error(it.message ?: "An unknown error occurred")
            }
        }
    }

    fun register(username: String, email: String, password: String) {
        viewModelScope.launch {
            _authState.value = AuthState.Loading
            val registerResult = authRepository.register(username, email, password)
            registerResult.onSuccess {
                // After successful registration, automatically log the user in
                login(email, password)
            }.onFailure {
                _authState.value = AuthState.Error(it.message ?: "An unknown error occurred")
            }
        }
    }

    fun resetState() {
        _authState.value = AuthState.Idle
    }
}

sealed class AuthState {
    object Idle : AuthState()
    object Loading : AuthState()
    object Success : AuthState()
    data class Error(val message: String) : AuthState()
}