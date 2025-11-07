package com.example.tempokotlin

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import com.example.tempokotlin.ui.theme.TempokotlinTheme

@Composable
fun SettingsScreen(modifier: Modifier = Modifier) {
    Box(
        modifier = modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text(text = "This is the Settings Page")
    }
}

@Preview(showBackground = true)
@Composable
fun SettingsScreenPreview() {
    TempokotlinTheme {
        SettingsScreen()
    }
}