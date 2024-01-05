#include "authentication-service.h"

void init_telco (void) __attribute__ ((constructor));

void init_telco(void) {
  telco_init ();
}

extern void * authenticate(void*,char*);

struct _GoAuthenticationService {
  GObject parent;
  void * callback;
  GThreadPool * pool;
};

GoAuthenticationService * telco_go_authentication_service_new (void * callback);
static void telco_go_authentication_service_iface_init (gpointer g_iface, gpointer iface_data);
static void telco_go_authentication_service_dispose (GObject * object);
static void telco_go_authentication_service_authenticate (GoAuthenticationService * service, const gchar * token,
    GCancellable * cancellable, GAsyncReadyCallback callback, gpointer user_data);
static gchar * telco_go_authentication_service_authenticate_finish (GoAuthenticationService * service, GAsyncResult * result,
    GError ** error);
static void telco_go_authentication_service_do_authenticate (GTask * task, GoAuthenticationService * self);

G_DEFINE_TYPE_EXTENDED (GoAuthenticationService, telco_go_authentication_service, G_TYPE_OBJECT, 0,
    G_IMPLEMENT_INTERFACE (TELCO_TYPE_AUTHENTICATION_SERVICE, telco_go_authentication_service_iface_init))


GoAuthenticationService * telco_go_authentication_service_new (void * callback) {
  GoAuthenticationService * service = NULL;

  service = g_object_new (TELCO_TYPE_GO_AUTHENTICATION_SERVICE, NULL);
  service->callback = callback;

  return service;
}

static void telco_go_authentication_service_iface_init (gpointer g_iface, gpointer iface_data){
  TelcoAuthenticationServiceIface * iface = g_iface;

  iface->authenticate = telco_go_authentication_service_authenticate;
  iface->authenticate_finish = telco_go_authentication_service_authenticate_finish;
}

static void
telco_go_authentication_service_class_init (GoAuthenticationServiceClass * klass)
{
  GObjectClass * object_class = G_OBJECT_CLASS (klass);

  object_class->dispose = telco_go_authentication_service_dispose;
}

static void telco_go_authentication_service_dispose (GObject * object) {
  GoAuthenticationService * self = TELCO_GO_AUTHENTICATION_SERVICE(object);
  
  if (self->pool != NULL) {
    g_thread_pool_free (self->pool, FALSE, FALSE);
    self->pool = NULL;
  }

  if (self->callback != NULL) {
    self->callback = NULL;
  }

  G_OBJECT_CLASS (telco_go_authentication_service_parent_class)->dispose (object);
}

static void
telco_go_authentication_service_init (GoAuthenticationService * self)
{
  self->pool = g_thread_pool_new ((GFunc) telco_go_authentication_service_do_authenticate, self, 1, FALSE, NULL);
}

static void telco_go_authentication_service_authenticate (GoAuthenticationService * service, const gchar * token, 
GCancellable * cancellable, GAsyncReadyCallback callback, gpointer user_data)
{
  GoAuthenticationService * self;
  GTask * task;

  self = TELCO_GO_AUTHENTICATION_SERVICE (service);

  task = g_task_new (self, cancellable, callback, user_data);
  g_task_set_task_data (task, g_strdup (token), g_free);

  g_thread_pool_push (self->pool, task, NULL);
}

static gchar *
telco_go_authentication_service_authenticate_finish (GoAuthenticationService * service, GAsyncResult * result, GError ** error)
{
  return g_task_propagate_pointer (G_TASK (result), error);
}

static void
telco_go_authentication_service_do_authenticate (GTask * task, GoAuthenticationService * self)
{
    const gchar * token;
    const gchar * session_info = NULL;
    gchar * message;
    void * result = NULL;

    token = g_task_get_task_data (task);

    result = authenticate(self->callback, (char*)token);

    if (result == NULL) {
        message = g_strdup ("Internal error");
    }

    session_info = (char*)result;
        
    if (session_info != NULL) {
        g_task_return_pointer (task, session_info, g_free);
    } else {
        g_task_return_new_error (task, TELCO_ERROR, TELCO_ERROR_INVALID_ARGUMENT, "%s", message);
    }

    g_free (message);
    g_object_unref (task);
}