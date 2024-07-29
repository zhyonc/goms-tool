#include "util.h"


void Log(const char* format, ...)
{
	char buf[1024] = { 0 };

	va_list args;
	va_start(args, format);
	vsprintf_s(buf, format, args);
	OutputDebugStringA(buf);
	//WriteText("ijl15.log", buf);
	va_end(args);
}

void WriteText(const char* filename, char* buf) {
	std::fstream file;
	file.open(filename, std::ios::out | std::ios::app);
	if (file.is_open()) {
		file  << buf << std::endl;
		file.close();
	}
}


BOOL SetHook(BOOL bInstall, PVOID* ppvTarget, PVOID pvDetour)
{
	if (DetourTransactionBegin() != NO_ERROR)
	{
		return FALSE;
	}

	auto tid = GetCurrentThread();

	if (DetourUpdateThread(tid) == NO_ERROR)
	{
		auto func = bInstall ? DetourAttach : DetourDetach;

		if (func(ppvTarget, pvDetour) == NO_ERROR)
		{
			if (DetourTransactionCommit() == NO_ERROR)
			{
				return TRUE;
			}
		}
	}
	DetourTransactionAbort();
	return FALSE;
}

DWORD GetFuncAddress(LPCSTR lpModule, LPCSTR lpFunc)
{
	auto mod = LoadLibraryA(lpModule);

	if (!mod)
	{
		return 0;
	}

	auto address = (DWORD)GetProcAddress(mod, lpFunc);

	if (!address)
	{
		return 0;
	}

#ifdef _DEBUG
	Log(__FUNCTION__" [%s] %s @ %8X", lpModule, lpFunc, address);
#endif
	return address;
}

#define JMP		0xE9
#define NOP		0x90
#define RET		0xC3
#define relative_address(frm, to) (int)(((int)to - (int)frm) - 5)

void PatchJmp(DWORD dwAddress, DWORD dwDest)
{

	*(BYTE*)dwAddress = JMP;
	*(DWORD*)(dwAddress + 1) = relative_address(dwAddress, dwDest);
}

void PatchRetZero(DWORD dwAddress)
{
	*(BYTE*)(dwAddress + 0) = 0x33;
	*(BYTE*)(dwAddress + 1) = 0xC0;
	*(BYTE*)(dwAddress + 2) = RET;
}

void PatchDMGCap(DWORD dwAddress)
{
	*(BYTE*)(dwAddress + 0) = 0x68;
	*(BYTE*)(dwAddress + 1) = 0xFF;
	*(BYTE*)(dwAddress + 2) = 0xFF;
	*(BYTE*)(dwAddress + 3) = 0xFF;
	*(BYTE*)(dwAddress + 4) = 0x7F;
}

void PatchChat(DWORD dwAddress, DWORD dwCount)
{
	for (int i = 0; i < dwCount; i++)
	{
		*(BYTE*)(dwAddress + i) = 0xEB;
	}
}

void PatchNop(DWORD dwAddress, DWORD dwCount)
{
	for (int i = 0; i < dwCount; i++)
	{
		*(BYTE*)(dwAddress + i) = NOP;
	}
}