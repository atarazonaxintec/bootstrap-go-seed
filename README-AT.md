# Segmento de cambios realizados
En este segmento se enumeraran los cambios, consideraciones, límites conocidos y evolutivos que puedan hacerse al proyecto; realizados por Alejandro Tarazona en conjunto con Copilot como asistente de desarrollo, y se le delegó la creación y edición de casos de prueba así como el completar líneas de código standard (dtos, entities y los useCases que fueran simples) ***Las decisiones arquitectónicas y procedimentales fueron enteramente realizadas por mi persona***.
## Cosas que hice
- Orders:
    - **Creación de orden vacía**: se crea inicialmente una orden sólo con el rut del cliente (que se usa como identificador del mismo) y que genera un UUID para futuras ediciones.
    - **Edición de orden**: acá se sobreecriben los productos asociados a la orden, siempre dejando la última actualización.
    - **Obtener órdenes por rut**: se listan las órdenes que tiene el cliente.
    - **Obtener orden por UUID**: se recupera la orden por su identificador.
- Payments:
    - **Creación de pago vacío**: se crea un pago vacío, asociado a una órden mediante su identificador.
    - **Pagar**: se registra la información de pago correspondiente
    - **Actualizar estado**: se actualiza el pago con la información recuperada de la respuesta del 'portal de pagos' (simulado).
- Products:
    - **Obtener todos los productos**: lista la información de los productos guardados.
    - **Obtener un producto por SKU**: muestra la informacón guardada de un solo producto.
- Pub/Subs:
    - Usé un emulador para simular la conexión a gcloud:
    ```shell
        brew install --cask google-cloud-sdk
        gcloud components install pubsub-emulator
        gcloud beta emulators pubsub start --project=demo-project --host-port=localhost:8085
    ```
- Ejecución
    La ejecución deberia hacerse usando el Makefile provisto y configurando el .env correspondiente según el **example.env**
- Tests
    Se deja una [Collection.json](./migrations/go-seed-orders.postman_collection.json) lista para hacer pruebas
## Consideraciones
1. Por parte de los productos, se hicieron los endpoints de recuperación para facilitar las pruebas.
2. Implementé una política de "creación vacía" + "edición/sobrrescritura" para manejar así la idempotencia y tolerar la concurrencia. Esto es: si un usuario reintenta la creación, se le asigna un nuevo UUID y si reintenta la edición queda con la última edición realizada.
## Límites conocidos
1. No supe dónde dejar los 'enums' que se usarían en las entidades y los reusé para los dto, lo cual no es elegante.
2. Está incompleto por muchas partes, falta integración con GCP y con algún portal de pagos.
3. No hice la integración con Pub/Subs, sólo emular la conexión para poder desplegar y hacer pruebas
## Siguientes pasos
1. Pruebas exahustivas a los flujos
2. Revisión de la comunicación con PubSubs
3. Simular el pago
4. Implementar el envío de la orden