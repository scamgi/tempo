package com.example.tempo

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountBox
import androidx.compose.material.icons.filled.DateRange
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.List
import androidx.compose.material.icons.filled.Menu
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.DrawerValue
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.ModalDrawerSheet
import androidx.compose.material3.ModalNavigationDrawer
import androidx.compose.material3.NavigationDrawerItem
import androidx.compose.material3.NavigationDrawerItemDefaults
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.adaptive.navigationsuite.NavigationSuiteScaffold
import androidx.compose.material3.rememberDrawerState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import androidx.compose.ui.unit.dp
import com.example.tempo.ui.theme.TempoTheme
import kotlinx.coroutines.launch

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            TempoTheme {
                TempoApp()
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@PreviewScreenSizes
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
                AppDestination.TODO -> TodoScreen(modifier = Modifier.padding(innerPadding))
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

enum class TodoDestination(
    val label: String,
    val icon: ImageVector,
) {
    HOME("Home", Icons.Default.Home),
    FAVORITES("Favorites", Icons.Default.Favorite),
    PROFILE("Profile", Icons.Default.AccountBox),
}

@Composable
fun TodoScreen(modifier: Modifier = Modifier) {
    var currentTodoDestination by rememberSaveable { mutableStateOf(TodoDestination.HOME) }

    NavigationSuiteScaffold(
        modifier = modifier,
        navigationSuiteItems = {
            TodoDestination.entries.forEach {
                item(
                    icon = {
                        Icon(
                            it.icon,
                            contentDescription = it.label
                        )
                    },
                    label = { Text(it.label) },
                    selected = it == currentTodoDestination,
                    onClick = { currentTodoDestination = it }
                )
            }
        }
    ) {
        Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
            when (currentTodoDestination) {
                TodoDestination.HOME -> HomeScreen(modifier = Modifier.padding(innerPadding))
                TodoDestination.FAVORITES -> FavoritesScreen(modifier = Modifier.padding(innerPadding))
                TodoDestination.PROFILE -> ProfileScreen(modifier = Modifier.padding(innerPadding))
            }
        }
    }
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

@Composable
fun HomeScreen(modifier: Modifier = Modifier) {
    Text(
        text = "Hello Home!",
        modifier = modifier
    )
}

@Composable
fun FavoritesScreen(modifier: Modifier = Modifier) {
    Text(
        text = "Hello Favorites!",
        modifier = modifier
    )
}

@Composable
fun ProfileScreen(modifier: Modifier = Modifier) {
    Text(
        text = "Hello Profile!",
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