package com.example.tempo.data

import com.example.tempo.data.remote.ApiService
import com.example.tempo.data.remote.LoginRequest
import com.example.tempo.data.remote.RegisterRequest

class AuthRepository(private val apiService: ApiService) {

    suspend fun login(email: String, password: String): Result<String> {
        return try {
            val response = apiService.login(LoginRequest(email, password))
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!.token)
            } else {
                Result.failure(Exception("Login failed: ${response.message()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun register(username: String, email: String, password: String): Result<Unit> {
        return try {
            val response = apiService.register(RegisterRequest(username, email, password))
            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(Exception("Registration failed: ${response.message()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}