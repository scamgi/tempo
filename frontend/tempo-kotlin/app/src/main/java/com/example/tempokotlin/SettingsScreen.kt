package com.example.tempokotlin

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.selection.selectable
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.RadioButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.example.tempokotlin.ui.theme.TempokotlinTheme
import com.example.tempokotlin.ui.theme.ThemeSetting
import java.util.Locale

@Composable
fun SettingsScreen(
    currentTheme: ThemeSetting,
    onThemeChange: (ThemeSetting) -> Unit,
    modifier: Modifier = Modifier
) {
    val themes = ThemeSetting.values()

    Column(modifier = modifier.padding(16.dp)) {
        Text("Theme", style = MaterialTheme.typography.titleLarge)
        Spacer(modifier = Modifier.height(8.dp))
        themes.forEach { theme ->
            Row(
                Modifier
                    .fillMaxWidth()
                    .selectable(
                        selected = (theme == currentTheme),
                        onClick = { onThemeChange(theme) }
                    )
                    .padding(vertical = 8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                RadioButton(
                    selected = (theme == currentTheme),
                    onClick = { onThemeChange(theme) }
                )
                Text(
                    text = theme.name.replaceFirstChar {
                        if (it.isLowerCase()) it.titlecase(
                            Locale.getDefault()
                        ) else it.toString()
                    },
                    style = MaterialTheme.typography.bodyLarge,
                    modifier = Modifier.padding(start = 16.dp)
                )
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun SettingsScreenPreview() {
    TempokotlinTheme {
        SettingsScreen(
            currentTheme = ThemeSetting.SYSTEM,
            onThemeChange = {}
        )
    }
}