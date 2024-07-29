#pragma once
#include "global.h"
#include "util.h"
#include <winsock2.h>
#include <ws2spi.h>
#include <ws2tcpip.h>

#pragma comment(lib, "ws2_32.lib")

namespace sock {
	bool HookWSPStartup(); 
	int WINAPI WSPConnect(SOCKET s, const struct sockaddr *name, int namelen, LPWSABUF lpCallerData, LPWSABUF lpCalleeData, LPQOS lpSQOS, LPQOS lpGQOS, LPINT lpErrno);
	int WINAPI WSPGetPeerName(SOCKET s, struct sockaddr *name, LPINT namelen, LPINT lpErrno);
}