package com.example.tempo.data

import com.example.tempo.data.remote.ApiService
import com.example.tempo.data.remote.CreateTodoListPayload
import com.example.tempo.data.remote.TodoList
import com.example.tempo.data.remote.TodoListWithItems
import kotlinx.coroutines.flow.first
import java.io.IOException

class TodoRepository(
    private val apiService: ApiService,
    private val userPreferencesRepository: UserPreferencesRepository
) {

    private suspend fun getAuthHeader(): String {
        val token = userPreferencesRepository.authToken.first()
        return "Bearer $token"
    }

    suspend fun getTodoLists(): Result<List<TodoList>> {
        return try {
            val response = apiService.getTodoLists(getAuthHeader())
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(IOException("Failed to fetch lists: ${response.message()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getTodoListWithItems(listId: Int): Result<TodoListWithItems> {
        return try {
            val response = apiService.getTodoListWithItems(getAuthHeader(), listId)
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(IOException("Failed to fetch list details: ${response.message()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun createTodoList(title: String): Result<TodoList> {
        return try {
            val response = apiService.createTodoList(getAuthHeader(), CreateTodoListPayload(title))
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(IOException("Failed to create list: ${response.errorBody()?.string()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}