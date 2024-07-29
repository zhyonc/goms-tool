// dllmain.cpp : Defines the entry point for the DLL application.
#include "ijl15.h"
#include "sock_hook.h"
#include "window_hook.h"

void WINAPI Patch(PVOID)
{
	if (!ijl15::CreateDllFile()) {
		Log("ijl15: Unable to write default dll to disk");
		return;
	}
	if (!ijl15::LoadDllFile()) {
		Log("ijl15: Unable to load default dll to memory");
		return;
	}
	if (!sock::HookWSPStartup()) { // Redirect the connect addr
		Log("Failed Hooking WSPStartup");
	}
	if (!window::HookCreateWindowExA()) { // WindowTitle
		Log("Failed Hooking window CreateWindowExA");
	}
	//if (!window::HookCreateMutexA()) {
	//	Log("Failed Hooking window CreateMutexA");// Multi-window conflict with skip sdo auth
	//}
}

BOOL APIENTRY DllMain(HMODULE hModule,
	DWORD  ul_reason_for_call,
	LPVOID lpReserved
)
{
	switch (ul_reason_for_call)
	{
	case DLL_PROCESS_ATTACH:
		DisableThreadLibraryCalls(hModule);
		CreateThread(NULL, NULL, (LPTHREAD_START_ROUTINE)&Patch, NULL, NULL, NULL);
		break;
	case DLL_THREAD_ATTACH:
	case DLL_THREAD_DETACH:
	case DLL_PROCESS_DETACH:
		break;
	}
	return TRUE;
}