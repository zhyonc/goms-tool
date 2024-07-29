#include "sock_hook.h"

namespace sock {

	WSPPROC_TABLE g_ProcTable;


	bool HookWSPStartup() {
		static auto _WSPStartup = decltype(&WSPStartup)(GetFuncAddress("MSWSOCK", "WSPStartup"));

		decltype(&WSPStartup) Hook = [](WORD wVersionRequested, LPWSPDATA lpWSPData, LPWSAPROTOCOL_INFOW lpProtocolInfo, WSPUPCALLTABLE UpcallTable, LPWSPPROC_TABLE lpProcTable) -> int
		{
			Log("[WSPStartup] Hijacked ProcTable");
			int ret = _WSPStartup(wVersionRequested, lpWSPData, lpProtocolInfo, UpcallTable, lpProcTable);
			g_ProcTable = *lpProcTable;

			lpProcTable->lpWSPConnect = WSPConnect;
			lpProcTable->lpWSPGetPeerName = WSPGetPeerName;

			return ret;
		};
		return SetHook(true, reinterpret_cast<void**>(&_WSPStartup), Hook);
	}

	int WINAPI WSPConnect(SOCKET s, const struct sockaddr *name, int namelen, LPWSABUF lpCallerData, LPWSABUF lpCalleeData, LPQOS lpSQOS, LPQOS lpGQOS, LPINT lpErrno)
	{
		sockaddr_in* service = (sockaddr_in*)name;
		char szAddr[INET_ADDRSTRLEN];
		inet_ntop(AF_INET, &(service->sin_addr), szAddr, INET_ADDRSTRLEN);

		if (strstr(szAddr, OPT_ADDR_SEARCH))
		{
			if (inet_pton(AF_INET, OPT_ADDR_HOSTNAME, &service->sin_addr) != 1) {
				Log("[WSPConnect] Invalid IP address");
			}
			else {
				Log("[WSPConnect] Replaced: %s", OPT_ADDR_HOSTNAME);
			}
		}
		else
		{
			Log("[WSPConnect] Original: %s", szAddr);
		}
		return g_ProcTable.lpWSPConnect(s, name, namelen, lpCallerData, lpCalleeData, lpSQOS, lpGQOS, lpErrno);
	}

	int WINAPI WSPGetPeerName(SOCKET s, struct sockaddr *name, LPINT namelen, LPINT lpErrno) {
		int nRet = g_ProcTable.lpWSPGetPeerName(s, name, namelen, lpErrno);

		if (nRet == SOCKET_ERROR)
		{
			Log("[WSPGetPeerName] ErrorCode: %d", *lpErrno);
		}
		else
		{
			sockaddr_in* service = (sockaddr_in*)name;
			char szAddr[INET_ADDRSTRLEN];
			inet_ntop(AF_INET, &(service->sin_addr), szAddr, INET_ADDRSTRLEN);

			auto nPort = ntohs(service->sin_port);

			if ( nPort >= OPT_PORT_LOW && nPort <= OPT_PORT_HIGH)
			{
				if (inet_pton(AF_INET, OPT_ADDR_NEXON, &service->sin_addr) != 1) {
					Log("[WSPGetPeerName] Invalid IP address");
				}
				else {
					Log("[WSPGetPeerName] Replaced: %s", OPT_ADDR_NEXON);
				}
			}
			else
			{
				Log("[WSPGetPeerName] Original: %s", szAddr);
			}
		}
		return  nRet;
	}
}