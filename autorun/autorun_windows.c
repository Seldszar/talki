#include <windows.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <objbase.h>
#include <shlobj.h>

char CreateShortcut(char *shortcutA, char *path, char *args) {
	IShellLink* pISL;
	IPersistFile* pIPF;
	HRESULT hr;

	CoInitializeEx(NULL, COINIT_MULTITHREADED);

	WORD shortcutW[MAX_PATH];
	int nChar = MultiByteToWideChar(CP_ACP, 0, shortcutA, -1, shortcutW, MAX_PATH);

	hr = CoCreateInstance(&CLSID_ShellLink, NULL, CLSCTX_INPROC_SERVER, &IID_IShellLink, (LPVOID*)&pISL);
	if (!SUCCEEDED(hr)) {
		return 1;
	}

	hr = pISL->lpVtbl->SetPath(pISL, path);
	if (!SUCCEEDED(hr)) {
		return 1;
	}

	hr = pISL->lpVtbl->SetArguments(pISL, args);
	if (!SUCCEEDED(hr)) {
		return 1;
	}

	hr = pISL->lpVtbl->QueryInterface(pISL, &IID_IPersistFile, (void**)&pIPF);
	if (!SUCCEEDED(hr)) {
		return 1;
	}

	hr = pIPF->lpVtbl->Save(pIPF, shortcutW, FALSE);
	if (!SUCCEEDED(hr)) {
		return 1;
	}

	pIPF->lpVtbl->Release(pIPF);
	pISL->lpVtbl->Release(pISL);

	return 0;
}
