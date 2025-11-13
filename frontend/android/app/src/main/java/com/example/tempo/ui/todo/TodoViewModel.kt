package com.example.tempo.ui.todo

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.tempo.data.TodoRepository
import com.example.tempo.data.remote.TodoList
import com.example.tempo.data.remote.TodoItem
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

// Represents the state of the Todo screen
data class TodoUiState(
    val isLoading: Boolean = false,
    val todoLists: List<TodoList> = emptyList(),
    val selectedListItems: List<TodoItem> = emptyList(),
    val selectedListTitle: String? = null,
    val errorMessage: String? = null,
    val selectedListId: Int? = null
)

class TodoViewModel(private val todoRepository: TodoRepository) : ViewModel() {

    private val _uiState = MutableStateFlow(TodoUiState())
    val uiState: StateFlow<TodoUiState> = _uiState.asStateFlow()

    init {
        fetchTodoLists()
    }

    fun fetchTodoLists() {
        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true, errorMessage = null) }
            val result = todoRepository.getTodoLists()
            result.onSuccess { lists ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        todoLists = lists
                    )
                }
                // Automatically select the first list if available
                if (lists.isNotEmpty() && _uiState.value.selectedListId == null) {
                    lists.firstOrNull()?.let { selectList(it.id) }
                } else if (lists.isEmpty()) {
                    // Clear selected list if no lists exist
                    _uiState.update { it.copy(selectedListId = null, selectedListItems = emptyList(), selectedListTitle = null) }
                }
            }.onFailure { e ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage = e.message ?: "An unknown error occurred"
                    )
                }
            }
        }
    }

    fun selectList(listId: Int) {
        viewModelScope.launch {
            // Avoid reloading if the list is already selected
            if (_uiState.value.selectedListId == listId) return@launch

            _uiState.update { it.copy(isLoading = true, errorMessage = null, selectedListId = listId) }
            val result = todoRepository.getTodoListWithItems(listId)
            result.onSuccess { listWithItems ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        selectedListItems = listWithItems.items,
                        selectedListTitle = listWithItems.title
                    )
                }
            }.onFailure { e ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage = e.message ?: "An unknown error occurred"
                    )
                }
            }
        }
    }

    fun createTodoList(title: String) {
        if (title.isBlank()) {
            _uiState.update { it.copy(errorMessage = "Title cannot be empty") }
            return
        }
        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true, errorMessage = null) }
            val result = todoRepository.createTodoList(title)
            result.onSuccess {
                // Refresh the lists to show the new one
                fetchTodoLists()
            }.onFailure { e ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage = e.message ?: "Failed to create list"
                    )
                }
            }
        }
    }
}