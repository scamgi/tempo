package com.example.tempo.data

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "user_prefs")

class UserPreferencesRepository(private val context: Context) {

    private object PreferencesKeys {
        val AUTH_TOKEN = stringPreferencesKey("auth_token")
    }

    val authToken: Flow<String?> = context.dataStore.data
        .map { preferences ->
            preferences[PreferencesKeys.AUTH_TOKEN]
        }

    suspend fun saveAuthToken(token: String) {
        context.dataStore.edit { preferences ->
            preferences[PreferencesKeys.AUTH_TOKEN] = token
        }
    }

    suspend fun clearAuthToken() {
        context.dataStore.edit { preferences ->
            preferences.remove(PreferencesKeys.AUTH_TOKEN)
        }
    }
}