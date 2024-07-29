#include "ijl15.h"

namespace ijl15 {
	DWORD ijlErrorStr_Proc;
	DWORD ijlFree_Proc;
	DWORD ijlGetLibVersion_Proc;
	DWORD ijlInit_Proc;
	DWORD ijlRead_Proc;
	DWORD ijlWrite_Proc;


	bool CreateDllFile() {
		DWORD dw;
		HANDLE hOrg = CreateFileA(LIB_NAME, (GENERIC_READ | GENERIC_WRITE), NULL, NULL, CREATE_ALWAYS, NULL, NULL);

		if (hOrg)
		{
			WriteFile(hOrg, g_Ijl15_Raw, sizeof(g_Ijl15_Raw), &dw, NULL);
			CloseHandle(hOrg);
			return true;
		}
		return false;
	}


	bool LoadDllFile() {

		HMODULE SeData_Base = LoadLibraryA(LIB_NAME);
		if (SeData_Base) {
			ijlErrorStr_Proc = (DWORD)GetProcAddress(SeData_Base, "ijlErrorStr");
			ijlFree_Proc = (DWORD)GetProcAddress(SeData_Base, "ijlFree");
			ijlGetLibVersion_Proc = (DWORD)GetProcAddress(SeData_Base, "ijlGetLibVersion");
			ijlInit_Proc = (DWORD)GetProcAddress(SeData_Base, "ijlInit");
			ijlRead_Proc = (DWORD)GetProcAddress(SeData_Base, "ijlRead");
			ijlWrite_Proc = (DWORD)GetProcAddress(SeData_Base, "ijlWrite");
			return true;
		}
		return false;
	}

#define LIB_EXPORT	extern "C" __declspec(dllexport)

	LIB_EXPORT void ijlGetLibVersion()
	{
		__asm 	 jmp dword ptr[ijlGetLibVersion_Proc]
	}

	LIB_EXPORT void ijlInit()
	{
		__asm  jmp dword ptr[ijlInit_Proc]
	}

	LIB_EXPORT void ijlFree()
	{
		__asm 	 jmp dword ptr[ijlFree_Proc]
	}

	LIB_EXPORT void ijlRead()
	{
		__asm jmp dword ptr[ijlRead_Proc]
	}

	LIB_EXPORT void ijlWrite()
	{
		__asm  jmp dword ptr[ijlWrite_Proc]
	}

	LIB_EXPORT void ijlErrorStr()
	{
		__asm  jmp dword ptr[ijlErrorStr_Proc]
	}

}
