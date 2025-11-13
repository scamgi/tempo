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
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import com.example.tempo.data.remote.TodoItem
import com.example.tempo.data.remote.TodoList

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TodoScreen(
    viewModel: TodoViewModel,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
    var showAddListDialog by remember { mutableStateOf(false) }

    if (showAddListDialog) {
        AddListDialog(
            onDismiss = { showAddListDialog = false },
            onConfirm = { title ->
                viewModel.createTodoList(title)
                showAddListDialog = false
            }
        )
    }

    Scaffold(
        modifier = modifier,
        floatingActionButton = {
            FloatingActionButton(onClick = { showAddListDialog = true }) {
                Icon(Icons.Default.Add, contentDescription = "Add List")
            }
        }
    ) { paddingValues ->
        Row(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            // Left Pane: List of Todo Lists
            Surface(
                modifier = Modifier.width(180.dp),
                color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.3f)
            ) {
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
                } else if (uiState.isLoading && uiState.todoLists.isEmpty()) {
                    Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                        CircularProgressIndicator()
                    }
                } else if (!uiState.isLoading && uiState.todoLists.isEmpty()) {
                    Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                        Text("Create your first list using the '+' button.", textAlign = TextAlign.Center)
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
}

@Composable
fun AddListDialog(
    onDismiss: () -> Unit,
    onConfirm: (String) -> Unit
) {
    var title by remember { mutableStateOf("") }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Create New List") },
        text = {
            OutlinedTextField(
                value = title,
                onValueChange = { title = it },
                label = { Text("List Title") },
                singleLine = true
            )
        },
        confirmButton = {
            Button(
                onClick = { onConfirm(title) },
                enabled = title.isNotBlank()
            ) {
                Text("Create")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("Cancel")
            }
        }
    )
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
    } else if (todoLists.isEmpty()) {
        Box(modifier = Modifier.fillMaxSize().padding(16.dp), contentAlignment = Alignment.Center) {
            Text(
                "No lists yet.",
                style = MaterialTheme.typography.bodyMedium,
                textAlign = TextAlign.Center,
                color = LocalContentColor.current.copy(alpha = 0.6f)
            )
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
                    color = if (list.id == selectedListId) MaterialTheme.colorScheme.primary else LocalContentColor.current,
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
    // This Scaffold is now inside the main screen's Scaffold,
    // which may not be ideal, but for this simple layout it works.
    // Consider restructuring if you add more complex features like a bottom bar here.
    Column(
        modifier = Modifier.fillMaxSize()
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
        } else if (items.isEmpty()) {
            Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                Text(
                    "No tasks in this list yet.",
                    style = MaterialTheme.typography.bodyMedium,
                    color = LocalContentColor.current.copy(alpha = 0.6f)
                )
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