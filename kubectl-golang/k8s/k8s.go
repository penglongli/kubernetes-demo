package k8s

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	base64KubeConfig = "YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VONVJFTkRRV0pEWjBGM1NVSkJaMGxDUVVSQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSYzBaQlJFRldUVkpOZDBWUldVUldVVkZFUlhkd2NtUlhTbXdLWTIwMWJHUkhWbnBOUWpSWVJGUkpkMDFFWTNkT2FrVjZUVVJqTVUweGIxaEVWRTEzVFVSamQwNUVSWHBOUkdNeFRURnZkMFpVUlZSTlFrVkhRVEZWUlFwQmVFMUxZVE5XYVZwWVNuVmFXRkpzWTNwRFEwRlRTWGRFVVZsS1MyOWFTV2gyWTA1QlVVVkNRbEZCUkdkblJWQkJSRU5EUVZGdlEyZG5SVUpCVFV0Q0NsQm1aMXBVYjFablJDdExaSE52VlZGeWMxcHNNQzlJZVZCVWJtc3ZXbFpFZW5CR05rSjNlSEZtZGtOT1FuSnVOM3BaVm5aYVMyNVFWRTlHYW01dmNFc0tNa2RhZVRKa2FWTkdhbTFYZFVWUlFVcDZhbFJCYTBGbVpVNUhhME5HZGpoc2IyNXVXV05tU1RJMVlXdHFWVXBtV2trMWMwczVjRmR4YldJMFdISTJZZ3BaY0V0bk5sQTRRbGhEU2podGIyRndjVzFHZVhKUVZVMUhMM0JHY0hvMGVUUkNRek15UkZGQk0xRXdia1ZNVjI5TFNrcFJTSFJuV2pWUlZtTkJXa3BWQ25SRGN6UjJkbmc1YUhaUFZXaDBUVzgxUlZJelRWQnZWRzQxYzFoMlJHMVFOV28xYUVaS1QweG9OelJZZGtoM2VIQnJjRnBoSzNkWVZHSXJZbFE1TWpZS1JVUmliak5IVTNwNlNWVlJXR3NyVmpkbE9GUkRTWGxRZVdGM04wdG9VMGRXVjBWaFdXazBhVmx2UzBKQldtWXJLMk5QYW0xNEwwOVdSREZYVFZGWWRRcG5VV05OWWtkT1IySnpSMEl6YlhkSWMxRTRRMEYzUlVGQllVMXFUVU5GZDBSbldVUldVakJRUVZGSUwwSkJVVVJCWjB0clRVRTRSMEV4VldSRmQwVkNDaTkzVVVaTlFVMUNRV1k0ZDBSUldVcExiMXBKYUhaalRrRlJSVXhDVVVGRVoyZEZRa0ZLTlhrclprMTZkV2xCY0dJeU9GbHphRE5rVG1VeFJFSkVVMFVLT1c5bGFUY3JTakU1TkU5VlUyNUhSamx3Y2tWSWFIWldRbEJNZEhCcE5tcFRkR0pVTDNwcE1XSlJWMlZ6UVhGTVpHMXhibkZGUWtRMGNFeFVVVnBUZFFwdVZ6azVkVTVXV1hOdlMzSmpVakJoYlVWelUzbFdWMGRqV1dZeGVXcFNZMkpaT0VGUlYwUlNkbTB4WjBRMU5YbHNla3RUTHpVd04zbGFVWEJ5VTFGYUNrNVpVMFEzZFdWNVYxYzFhSE5FVEc1WWJXaG1WamxsV0VsNFZFbFNTM2RDZDBOWlRFazBiRGM1UVZCRmVsaHBWMmwzYUhWeFpHdFhSR1pvVkhsVmFrMEtRVEJqY21aMFRIcFBSRTFDVUcwMlVtUXJha1prY21aNFRrWmlRbE4zYW1kUldtUTNkbWx6UjNKS2FXczRXamxvWW1Fdk1GY3hZbmgyYjNwb1dXWlZad280UTJNekwxUTJiazVNYzFkMGFIcHJaMEV2U0dGTFVIUndRbXRhWlhsaWNsRk5RWFUzYlhRNUwzbHhhM05YTkRORldqbDBiaTgzV0VRd1RUMEtMUzB0TFMxRlRrUWdRMFZTVkVsR1NVTkJWRVV0TFMwdExRbz0KICAgIHNlcnZlcjogaHR0cHM6Ly8xNzIuMjAuMTAuMTQ6NjQ0MwogIG5hbWU6IGt1YmVybmV0ZXMKY29udGV4dHM6Ci0gY29udGV4dDoKICAgIGNsdXN0ZXI6IGt1YmVybmV0ZXMKICAgIHVzZXI6IGt1YmVybmV0ZXMtYWRtaW4KICBuYW1lOiBrdWJlcm5ldGVzLWFkbWluQGt1YmVybmV0ZXMKY3VycmVudC1jb250ZXh0OiBrdWJlcm5ldGVzLWFkbWluQGt1YmVybmV0ZXMKa2luZDogQ29uZmlnCnByZWZlcmVuY2VzOiB7fQp1c2VyczoKLSBuYW1lOiBrdWJlcm5ldGVzLWFkbWluCiAgdXNlcjoKICAgIGNsaWVudC1jZXJ0aWZpY2F0ZS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VNNGFrTkRRV1J4WjBGM1NVSkJaMGxKVUVNNVRXTXhiVTlzVTFGM1JGRlpTa3R2V2tsb2RtTk9RVkZGVEVKUlFYZEdWRVZVVFVKRlIwRXhWVVVLUVhoTlMyRXpWbWxhV0VwMVdsaFNiR042UVdWR2R6QjVUVVJCTTAxRVdYaE5la0V6VGxST1lVWjNNSGxOVkVFelRVUlplRTE2UVROT1ZHUmhUVVJSZUFwR2VrRldRbWRPVmtKQmIxUkViazQxWXpOU2JHSlVjSFJaV0U0d1dsaEtlazFTYTNkR2QxbEVWbEZSUkVWNFFuSmtWMHBzWTIwMWJHUkhWbnBNVjBackNtSlhiSFZOU1VsQ1NXcEJUa0puYTNGb2EybEhPWGN3UWtGUlJVWkJRVTlEUVZFNFFVMUpTVUpEWjB0RFFWRkZRVEE0VkRCeWVYQkRhMkpqVEVNMU9IVUtPVGRxYVRZMmEzQkZaWEJXTDA5MVUydFlLMEp0V1hZNGMyVlBRMVZqTlM5R1RUVXpPSFYxTHpGWVRGTXlPVUpvZGk5YVJ6UkNVSEJxY1V3d1Z5OVFWUXA2VGpsWFdXVnFSVXRhYUdaNVdsQlVOM0Y1VWxwVlVWcFdNVWxWYmxreFpuZzViMUpUYzJsSVlqSXZWakpLUVVOV2IydEpaV1J4Ym1kUU5qUTRPVzh3Q2xnMFJreFVXbGhWZDFBMlpGSnlVV3RMYUZoWlowOUZjRXRTTjJaUGJtdzJZVk5ZTkhZM1MyOHdRVlJCUVc1RGRuUnlhWEZuZWtreE9IaE9TVkJaY2xVS1RWWlJkMkpqWkhFeWNGWjFibHB5U2pGUVVDdDRjQ3RvWVZKc2JUbFVZVWhqT0hsVlluTjBRMFk0UTBWeVdXRnlVeTh2ZVV0cVEzTnZRVElyVmtkMVl3cDFMemxyVnpaeE4wSkJiWFZEV0NzM1IwZG9OVFpQWkdaMk1DdDRjRzFxTkdsNmQwTm5PSEozTVcxbU9XbGFUbTVsVGpaVFdtTnpMMnAzT0V0bmVEaDZDamhXZW1RelVVbEVRVkZCUW05NVkzZEtWRUZQUW1kT1ZraFJPRUpCWmpoRlFrRk5RMEpoUVhkRmQxbEVWbEl3YkVKQmQzZERaMWxKUzNkWlFrSlJWVWdLUVhkSmQwUlJXVXBMYjFwSmFIWmpUa0ZSUlV4Q1VVRkVaMmRGUWtGQk4xaEZRMFV4VEdWcWVYSjRaelZ6Tm5sbVlrczNVV0kyVDFaWE1tOURiV3BrZVFwelNUbEdaRVJpYkZCRlMyVnNialJQTVM5TVdFNDVZM3BQY25FM09GaFBRbkpSWnpSTUswdDNhbnBoZERSeWFYRllVVU5IYWtwdVMxaEtXRzU1Y0dkUUNsa3pjeXROTHpkR2NtdHVjbGw2Y3poU1VXZE9abWxHVkRSVlNERm1ZM2c1VG5Bck0weFFUVmhqVkZKTmNGZHlVekZhZGxwbFNYaEdjWFkyTmpsU1RXb0tLM3BtTVRseWFFOW9VMGg1TjNKd1lteFlRWGMyYXpodFJGTnVSbHBIWW5GWlozSXpWRUV3Y0dObGEwVlBSSEl6WjNJMVlURk1WalpuU0hocFRGTjNhUW8xZVdFeGNUZHRRMkoyUkVoNFlsTnFWMGx6YldsV2FFRk1NbU42UW5CbGJtY3ZNbmRuYWxGd2Rqa3JXbUo2Wmt4TmRXRkZlRGtyWlRacFMyRkRkUzh5Q21SRU0yWXZXRzFqYnpOek5HbENVRzlXZFVkTVVGUlpTVWh3UVN0YVkzSXhibXhXTVVGdWVVOVJSbTEzTVVOUmFtUXlZejBLTFMwdExTMUZUa1FnUTBWU1ZFbEdTVU5CVkVVdExTMHRMUW89CiAgICBjbGllbnQta2V5LWRhdGE6IExTMHRMUzFDUlVkSlRpQlNVMEVnVUZKSlZrRlVSU0JMUlZrdExTMHRMUXBOU1VsRmNFRkpRa0ZCUzBOQlVVVkJNRGhVTUhKNWNFTnJZbU5NUXpVNGRUazNhbWsyTm10d1JXVndWaTlQZFZOcldDdENiVmwyT0hObFQwTlZZelV2Q2taTk5UTTRkWFV2TVZoTVV6STVRbWgyTDFwSE5FSlFjR3B4VERCWEwxQlZlazQ1VjFsbGFrVkxXbWhtZVZwUVZEZHhlVkphVlZGYVZqRkpWVzVaTVdZS2VEbHZVbE56YVVoaU1pOVdNa3BCUTFadmEwbGxaSEZ1WjFBMk5EZzViekJZTkVaTVZGcFlWWGRRTm1SU2NsRnJTMmhZV1dkUFJYQkxVamRtVDI1c05ncGhVMWcwZGpkTGJ6QkJWRUZCYmtOMmRISnBjV2Q2U1RFNGVFNUpVRmx5VlUxV1VYZGlZMlJ4TW5CV2RXNWFja294VUZBcmVIQXJhR0ZTYkcwNVZHRklDbU00ZVZWaWMzUkRSamhEUlhKWllYSlRMeTk1UzJwRGMyOUJNaXRXUjNWamRTODVhMWMyY1RkQ1FXMTFRMWdyTjBkSGFEVTJUMlJtZGpBcmVIQnRhalFLYVhwM1EyYzRjbmN4YldZNWFWcE9ibVZPTmxOYVkzTXZhbmM0UzJkNE9IbzRWbnBrTTFGSlJFRlJRVUpCYjBsQ1FWRkRlbFJGTjJaQlEycGpkSFF6YWdwUFUxQk1SMlk0U0VOSVNqbGxTM0pXVDFvNFprVmFXSEJMTVRCSlVVWm9WMkY2SzNSdWFVcDNlWEZ1YUZSNFlUUm9abGs1VlZsamQzTmhkRTR5VTNGTUNuTkRZVGhVTVhkVlEyTkJUV1EzWVdsT1ZXUTNSRk5GVGxoR2MxbFZObUZuZG5SSldtODRhVUZUVFdocE1Ga3hPV3B1UW5OV1FtMWFSV3RvUVV3eWNuQUtSRUppVEhobGVUaEJPUzlsYURkVmFEVXJUekoyYTJoUk1XeE1SRkJLV0VkMVpHWm9ObFphUkVkeVZHOW1UbE15WnpKRmRIUkxTbUpXUTJVeFlrWmxUd28xVWtRNGFraEpNbmw0U1hGRk9VUXJZblZYYW1wSFVpdFFlbVYyVERSQk4xRkZRM0poWldFMWRtMUhTMDV3WmxNMWIzcElSVGh3WWtwaWNFOUtZVUpFQ2xGVFRWZGhVRFp6YkRCaVEwcGxjM1EyWjNFM0wxRkpZaXN6YzNSS0wzVk5WM1JQWlRjcmNtSlpaVkZ2ZEdWdWRrOVNWblZLWkhOckx6QXdORkZzTUdFS1NVSjVNR05YUVVKQmIwZENRVTlHVEhCUlJHUlFlRzlYTUVGaGJYWkpUM0V5TkdwWGEycDNLM2RyV1d0cGRWbEZjbGhaWWs0eVJuQnRZbTVuUTFCV1N3b3lPRFJSYTFkRWVXUmhjRGN5Y0RjclRYcERURnBzWlU5UGRXMTZURGM0VTBoRVUzQlhlRUZhVFRVM1dGZHJVbkZ2TlRoNVUweDFNMlZRUVdwQ2JYbFhDa2xsTUVwTWNqVk5PRzlZTDBSaU9DdEVkbTVDUVRBM1RGRXZRa2R0Y0ZsQ1RrWTNjMlU0V21ORVJtNUpjRWwwTnk4d2JXbFdXRTQ1UVc5SFFrRlFRMmdLV2xwMFpqVndlRkZZTlVodlJqbE1PVWhzWmpGdGNrZEVPWFI1VWpaUFEyVjJOMFlyVGxKdE1HUTRUbFkyUms1NGRpOHJSbEYyZUVWTk9FWkZkMjk0Y2dwRVYxaEVORzVoVmlzdk5VeHdjMHRqVkhsaUwwNVllRWQ2VWpjelRrOU1XRlZuUVZkT1R5OHZWbEI2TVV4cFpqTlJWRXRUVTJ0aVkyOWpXRWg0WWpkS0NpczJaME5VYjBoRmEzaDRSWEoxWms0MVRWVk1jVEpSVG1Sc1RWRmtVV05GVERBdlNtZEhTR2hCYjBkQlREbHNaMEozYmpKWUwxVmxXVE4wUVV0SWNTc0tUelZpVHpScGNUSkhhVEpwY20weWNEVmFObk5MVnpGRVRDODBUVU5SVkVsSlJVUlhhalFyVUZCQ1JXOUJiMHRqYm5BdlRYTmhXRTFxTVhZeVRsRXJSUXB1UVRadUx6Um9TM1JXTm5Cc2ExSk1NR2RFWXk4M2JHTXplWFZvVEhGcVNVNWpObmRrWTNocU5tUXJZM0pPVW1sWE9YTktaRGdyYkRoNk9HSnFLMXBpQ20xQ2JVUjRSVXBEWjFWTmFUVXhhMFJvUlVWT1FYZFZRMmRaUlVFeGVua3lkM2hWTUVaUGRqZHZSM0JJVEZNelJpODBNbkk0YXpoRkswZFVNMDVxZDNBS1dWWXhkbkFyT1RaR1RuQTFPVzVJVUN0SlNUa3lZVWw2TTFKT2FFcG1URE40SzJoQldVTjRiV1pDT0RGTFpYaG9SRXBRU2xKelIxbzFja2hsZERJd09BcGFkVVl5VG5Cd1oySmxablJuYTJkNlNrVlpaMlJIVVZSMUwwdHRiV0pNV1VSSU9FeDBRakpsVEZCSWRYUldMMEV5ZEZKSFYxUmxSbVZIVmtoSFMwTlJDbkJWUTNKUGJVVkRaMWxCTUdZMU5HUmxXRnBXUWxWRmFHY3JjbGd6V2xKMVJIUXZLMjF0TTBOcGFrOUlWVFZhYnpRck5WSlVibk5rVjFvd1pFUlRka0VLWjFaWE1WQjVlbFpIYjAxT09GTXlUV0ZVZVN0Slp5OHJZbUVyVWt0SlVWQTJOWGRMYldnd2NFRmphVkZOYUdGMWNGRlFiWFZSYkdSUlpHZFlVMGRuUkFwaFlUaFNUSGtyWm01bmNUZDJhbXM1V2paRVdVdENURTg1T1c5c1ZIWjVXSFZFTWxsVVREUkhSR1JIUnpoTmNHeHFhMkpxTVVFOVBRb3RMUzB0TFVWT1JDQlNVMEVnVUZKSlZrRlVSU0JMUlZrdExTMHRMUW89"

	discoveryMapper *restmapper.DeferredDiscoveryRESTMapper
	dynamicClient dynamic.Interface
)

func GetDynamicClient() (dynamic.Interface, error) {
	if dynamicClient != nil {
		return dynamicClient, nil
	}

	restConfig, err := buildRestConfig(base64KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "build restConfig failed")
	}

	return dynamic.NewForConfig(restConfig)
}

func GetDiscoveryMapper() (*restmapper.DeferredDiscoveryRESTMapper, error) {
	if discoveryMapper != nil {
		return discoveryMapper, nil
	}

	restConfig, err := buildRestConfig(base64KubeConfig)
	if err != nil {
		return discoveryMapper, errors.Wrap(err, "build restConfig failed")
	}

	// Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return discoveryMapper, errors.Wrap(err, "new dc failed")
	}

	discoveryMapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	return discoveryMapper, nil
}

func buildRestConfig(base64KubeConfig string) (resetConfig *rest.Config, err error) {
	kubeConfig, err := base64.StdEncoding.DecodeString(base64KubeConfig)
	if err != nil {
		return nil, err
	}

	conf, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (config *clientcmdapi.Config, e error) {
		return clientcmd.Load(kubeConfig)
	})

	if err != nil {
		return nil, err
	}
	return conf, nil
}
