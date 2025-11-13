package com.example.tempo.data.remote

import com.google.gson.annotations.SerializedName
import java.util.Date

// Represents a single to-do list
data class TodoList(
    val id: Int,
    @SerializedName("userId")
    val userId: Int,
    val title: String,
    @SerializedName("createdAt")
    val createdAt: Date
)

// Represents a single item within a to-do list
data class TodoItem(
    val id: Int,
    @SerializedName("listId")
    val listId: Int,
    val task: String,
    @SerializedName("isCompleted")
    val isCompleted: Boolean,
    @SerializedName("dueDate")
    val dueDate: Date?,
    val priority: Int,
    @SerializedName("createdAt")
    val createdAt: Date
)

// Represents a to-do list along with its items, as returned by the /api/lists/{listId} endpoint
data class TodoListWithItems(
    val id: Int,
    @SerializedName("userId")
    val userId: Int,
    val title: String,
    @SerializedName("createdAt")
    val createdAt: Date,
    val items: List<TodoItem>
)
