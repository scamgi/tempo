package com.example.tempo.data.remote

import com.google.gson.annotations.SerializedName
import retrofit2.Response
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Path

// Base URL must point to the host machine's localhost from the Android emulator
private const val BASE_URL = "http://10.0.2.2:8080/api/"

interface ApiService {
    @POST("users/register")
    suspend fun register(@Body registerRequest: RegisterRequest): Response<Unit>

    @POST("users/login")
    suspend fun login(@Body loginRequest: LoginRequest): Response<LoginResponse>

    // --- To-Do List Endpoints ---

    @GET("lists")
    suspend fun getTodoLists(@Header("Authorization") token: String): Response<List<TodoList>>

    @GET("lists/{listId}")
    suspend fun getTodoListWithItems(
        @Header("Authorization") token: String,
        @Path("listId") listId: Int
    ): Response<TodoListWithItems>

    @POST("lists")
    suspend fun createTodoList(
        @Header("Authorization") token: String,
        @Body payload: CreateTodoListPayload
    ): Response<TodoList>
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

data class CreateTodoListPayload(
    val title: String
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