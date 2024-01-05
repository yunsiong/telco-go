#ifndef __AUTHENTICATION_SERVICE_H__
#define __AUTHENTICATION_SERVICE_H__

#include <telco-core.h>

#define TELCO_TYPE_GO_AUTHENTICATION_SERVICE (telco_go_authentication_service_get_type ())
G_DECLARE_FINAL_TYPE (GoAuthenticationService, telco_go_authentication_service, TELCO, GO_AUTHENTICATION_SERVICE, GObject)

GoAuthenticationService * telco_go_authentication_service_new (void * callback);

#endif 