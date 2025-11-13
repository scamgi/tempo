package com.example.tempo.data.remote

import com.google.gson.annotations.SerializedName
import retrofit2.Response
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.Body
import retrofit2.http.POST

// Base URL must point to the host machine's localhost from the Android emulator
private const val BASE_URL = "http://10.0.2.2:8080/api/"

interface ApiService {
    @POST("users/register")
    suspend fun register(@Body registerRequest: RegisterRequest): Response<Unit>

    @POST("users/login")
    suspend fun login(@Body loginRequest: LoginRequest): Response<LoginResponse>
}

// --- Data Transfer Objects ---

data class RegisterRequest(
    val username: String,
    val email: String,
    val password: String
)

data class LoginRequest(
    val email: String,
    val password: String
)

data class LoginResponse(
    @SerializedName("token")
    val token: String
)


object RetrofitInstance {
    val api: ApiService by lazy {
        Retrofit.Builder()
            .baseUrl(BASE_URL)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(ApiService::class.java)
    }
}