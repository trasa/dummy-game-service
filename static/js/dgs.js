;(function($) {


    $.logviewer = function (element, options) {

        var defaults = {
            foo: 'bar',
            socketUrl: "ws://" + window.location.host + "/viewlog",
            pingUrl: "http://" + window.location.host + "/ping",

            // if your plugin is event-driven, you may provide callback capabilities
            // for its events. execute these functions before or after events of your
            // plugin, so that users may customize those particular events without
            // changing the plugin's code
            onFoo: function () {
            }


        };

        // to avoid confusions, use "plugin" to reference the
        // current instance of the object
        var plugin = this;

        // this will hold the merged default, and user-provided options
        // plugin's properties will be available through this object like:
        // plugin.settings.propertyName from inside the plugin or
        // element.data('pluginName').settings.propertyName from outside the plugin,
        // where "element" is the element the plugin is attached to;
        plugin.settings = {};


        var $element = $(element), // reference to the jQuery version of DOM element
            element = element;    // reference to the actual DOM element

        var websocket;

        // ctor
        plugin.init = function () {
            plugin.settings = $.extend({}, defaults, options);
            console.log("init");
        };

        /* private */
        var displayMessage = function (msg) {
            plugin.settings["displayMessage"](msg);
        }

        /* public */
        plugin.run = function () {
            displayMessage("Connecting to " + plugin.settings["socketUrl"]);
            websocket = new WebSocket(plugin.settings["socketUrl"]);

            websocket.onopen = function (evt) {
                displayMessage("opening socket...ready");
            };
            websocket.onclose = function (evt) {
                displayMessage("Socket closed");
                websocket = null;
            };
            websocket.onmessage = function (evt) {
                console.log("msg rec: ");
                console.log(evt.data);
                displayMessage(evt.data)
            };
            websocket.onerror = function (evt) {
                displayMessage("Error Received");
                console.log("Error:");
                console.log(evt);
            };
        };

        /* public */
        plugin.send = function(msg) {
            msg = JSON.stringify(msg);
            console.log("send " + msg);
            websocket.send(msg);
        };

        /* public */
        plugin.forceClose = function() {
            console.log("force close!");
            websocket.close();
        };

        plugin.ping = function() {
            console.log("ping");
            $.get(plugin.settings["pingUrl"], function() {
                console.log("ping sent");
            })
        };

        plugin.init();
    };


    $.fn.logviewer = function (options) {
        return this.each(function () {
            if (undefined == $(this).data('logviewer')) {
                var plugin = new $.logviewer(this, options);
                $(this).data('logviewer', plugin);
            }
        });
    };

    var app = $.sammy(function () {
        this.use(Sammy.EJS);

        this.get('#/', function () {
            this.render('templates/viewlog.ejs', function (html) {
                $('#mainContent').html(html);
                var mc;

                $('#connect').click(function () {
                    mc = $('#connect').logviewer({
                        displayMessage: displayMessage
                    }).data('logviewer');
                    mc.run();
                });

                $('#disconnect').click(function() {
                    mc.forceClose();
                });

                $('#ping').click(function() {
                    mc.ping();
                });

                $('#calculateHash').click(function() {
                    var shaObj = new jsSHA("SHA-256", "TEXT");
                    shaObj.update($('#secretKey').val());
                    shaObj.update($('#hook_body').val());
                    displayMessage("hash is '" + shaObj.getHash("HEX") + "'");
                });

                function displayMessage(msg) {
                    var display = $('#display');
                    display
                        .append(msg + "\n\n")
                        .stop()
                        .animate({scrollTop: display[0].scrollHeight}, 800);
                }
            });
        });
    });

    $(function () {
        app.run('#/');
    });
})(jQuery);