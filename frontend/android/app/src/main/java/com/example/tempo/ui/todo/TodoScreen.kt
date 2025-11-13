package com.example.tempo.ui.todo

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.outlined.Check
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import com.example.tempo.data.remote.TodoItem
import com.example.tempo.data.remote.TodoList

@Composable
fun TodoScreen(
    viewModel: TodoViewModel,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()

    Row(modifier = modifier.fillMaxSize()) {
        // Left Pane: List of Todo Lists
        Box(modifier = Modifier.width(150.dp)) {
            TodoListPane(
                todoLists = uiState.todoLists,
                selectedListId = uiState.selectedListId,
                onSelectList = viewModel::selectList,
                isLoading = uiState.isLoading && uiState.todoLists.isEmpty()
            )
        }


        // Right Pane: Items of the selected list
        Column(modifier = Modifier
            .fillMaxHeight()
            .weight(1f)
            .padding(horizontal = 16.dp)) {
            if (uiState.selectedListTitle != null) {
                TodoItemsPane(
                    title = uiState.selectedListTitle!!,
                    items = uiState.selectedListItems,
                    isLoading = uiState.isLoading && uiState.selectedListItems.isEmpty()
                )
            } else if (uiState.isLoading) {
                Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator()
                }
            }
        }
    }

    if (uiState.errorMessage != null) {
        // Using a Box to show the error message at the bottom
        Box(modifier = Modifier.fillMaxSize().padding(16.dp), contentAlignment = Alignment.BottomCenter) {
            Text(
                text = uiState.errorMessage!!,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.error
            )
        }
    }
}

@Composable
fun TodoListPane(
    todoLists: List<TodoList>,
    selectedListId: Int?,
    onSelectList: (Int) -> Unit,
    isLoading: Boolean
) {
    if (isLoading) {
        Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            CircularProgressIndicator()
        }
    } else {
        LazyColumn(
            modifier = Modifier
                .fillMaxHeight()
                .padding(top = 8.dp)
        ) {
            items(todoLists) { list ->
                Text(
                    text = list.title,
                    fontWeight = if (list.id == selectedListId) FontWeight.Bold else FontWeight.Normal,
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { onSelectList(list.id) }
                        .padding(16.dp)
                )
            }
        }
    }
}

@Composable
fun TodoItemsPane(
    title: String,
    items: List<TodoItem>,
    isLoading: Boolean,
) {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(onClick = { /* TODO: Add new item functionality */ }) {
                Icon(Icons.Default.Add, contentDescription = "Add Task")
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .padding(paddingValues)
                .fillMaxSize()
        ) {
            Text(
                text = title,
                style = MaterialTheme.typography.headlineMedium,
                modifier = Modifier.padding(bottom = 16.dp, top = 8.dp)
            )

            if (isLoading) {
                Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator()
                }
            } else {
                LazyColumn {
                    items(items) { item ->
                        TodoItemRow(item = item)
                        Divider()
                    }
                }
            }
        }
    }
}

@Composable
fun TodoItemRow(item: TodoItem) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            imageVector = if (item.isCompleted) Icons.Filled.CheckCircle else Icons.Outlined.Check,
            contentDescription = "Task status",
            tint = if (item.isCompleted) MaterialTheme.colorScheme.primary else LocalContentColor.current,
            modifier = Modifier.size(24.dp)
        )
        Spacer(modifier = Modifier.width(16.dp))
        Text(
            text = item.task,
            textDecoration = if (item.isCompleted) TextDecoration.LineThrough else null,
            color = if (item.isCompleted) Color.Gray else LocalContentColor.current
        )
    }
}