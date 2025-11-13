package com.example.tempo

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.example.tempo.navigation.AppNavigation
import com.example.tempo.navigation.getTodoViewModel
import com.example.tempo.ui.theme.TempoTheme
import com.example.tempo.ui.todo.TodoScreen
import kotlinx.coroutines.launch

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            TempoTheme {
                // The AppNavigation composable now handles everything
                AppNavigation()
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TempoApp() {
    val drawerState = rememberDrawerState(initialValue = DrawerValue.Closed)
    val scope = rememberCoroutineScope()
    var currentDestination by rememberSaveable { mutableStateOf(AppDestination.TODO) }
    val currentTitle by remember(currentDestination) {
        derivedStateOf { currentDestination.label }
    }

    ModalNavigationDrawer(
        drawerContent = {
            ModalDrawerSheet {
                Spacer(Modifier.height(12.dp))
                AppDestination.entries.forEach { destination ->
                    NavigationDrawerItem(
                        icon = { Icon(destination.icon, contentDescription = null) },
                        label = { Text(destination.label) },
                        selected = destination == currentDestination,
                        onClick = {
                            scope.launch { drawerState.close() }
                            currentDestination = destination
                        },
                        modifier = Modifier.padding(NavigationDrawerItemDefaults.ItemPadding)
                    )
                }
            }
        },
        drawerState = drawerState
    ) {
        Scaffold(
            topBar = {
                TopAppBar(
                    title = { Text(currentTitle) },
                    navigationIcon = {
                        IconButton(onClick = { scope.launch { drawerState.open() } }) {
                            Icon(
                                imageVector = Icons.Default.Menu,
                                contentDescription = "Menu"
                            )
                        }
                    }
                )
            }
        ) { innerPadding ->
            when (currentDestination) {
                AppDestination.TODO -> TodoScreen(viewModel = getTodoViewModel(), modifier = Modifier.padding(innerPadding))
                AppDestination.NOTES -> NotesScreen(modifier = Modifier.padding(innerPadding))
                AppDestination.JOURNAL -> JournalScreen(modifier = Modifier.padding(innerPadding))
                AppDestination.SETTINGS -> SettingsScreen(modifier = Modifier.padding(innerPadding))
            }
        }
    }
}


enum class AppDestination(
    val label: String,
    val icon: ImageVector,
) {
    TODO("Todo", Icons.Filled.List),
    NOTES("Notes", Icons.Filled.Edit),
    JOURNAL("Journal", Icons.Filled.DateRange),
    SETTINGS("Settings", Icons.Filled.Settings),
}

@Composable
fun NotesScreen(modifier: Modifier = Modifier) {
    Text(
        text = "Hello Notes!",
        modifier = modifier
    )
}

@Composable
fun JournalScreen(modifier: Modifier = Modifier) {
    Text(
        text = "Hello Journal!",
        modifier = modifier
    )
}

@Composable
fun SettingsScreen(modifier: Modifier = Modifier) {
    Text(
        text = "Hello Settings!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    TempoTheme {
        TempoApp()
    }
}