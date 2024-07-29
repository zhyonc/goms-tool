#pragma once
#define WIN32_LEAN_AND_MEAN
#include <Windows.h>
#include <fstream>
#include "detours.h"
#pragma comment(lib, "detours.lib")


void Log(const char* format, ...);
void WriteText(const char* filename, char* buf);
BOOL SetHook(BOOL bInstall, PVOID* ppvTarget, PVOID pvDetour);
DWORD GetFuncAddress(LPCSTR lpModule, LPCSTR lpFunc);
void PatchJmp(DWORD dwAddress, DWORD dwDest);
void PatchRetZero(DWORD dwAddress);
void PatchDMGCap(DWORD dwAddress);
void PatchChat(DWORD dwAddress, DWORD dwCount);
void PatchNop(DWORD dwAddress, DWORD dwCount);