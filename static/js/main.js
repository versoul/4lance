
var lastNotifySoundTime = Date.now();
var notifySound = false;

$( document ).ready(function() {
    var dropRight = true;


    if($(window).width() < 977){
        dropRight = false;
    }
    var sound = new Howl({
        src: ['/static/media/5.mp3'],
        volume: 0.5,
    });

    $('.categoriesMultiselect').multiselect({
        enableClickableOptGroups: true,
        includeSelectAllOption: true,
        enableFiltering: true,
        buttonWidth: '100%',
        maxHeight: 400,
        dropRight: dropRight,
        numberDisplayed: 0,
        selectAllText: 'Выбрать все категории',
        nonSelectedText: 'Выбрать категории',
        nSelectedText: ' выбрано',
        allSelectedText: 'Выбраны все категории',
        filterPlaceholder: 'Поиск',
        filterBehavior: 'text',
        enableCaseInsensitiveFiltering: true,
        numberDisplayed: 0,
        buttonTitle: function(options, select){
            return "";
        }
    });

    if(!window.auth){
        $( "div.projectsFilter" ).tooltip({
            track:true
        });
    }
    else{
        $(".projectsFilter").prop("title", "");
    }


    $("#projectsList").click(function(e){
        if(!e.ctrlKey){
            e.preventDefault();
            var row = $(e.target).closest("tr");
            var lnk = row.first().first().find("a");
            var modal = $("#projectModal");
            modal.find(".modal-title").html(lnk.html());
            modal.find(".modal-body").html(row.data("description"));
            modal.find("#toProjectBtn").attr("href", lnk.attr("href"));
            modal.modal('show');
        }
    });
    $("#filterSaveBtn").click(function(e){
        e.preventDefault();
        var keyWords = $("#keyWords").tagsinput('items');
        var categories = [];
        var multiselects = $('.categoriesMultiselect');
        for(var i=0,l=multiselects.length; i<l; i++){
            categories = categories.concat($(multiselects[i]).val());
        }
        $.ajax({
            method: "POST",
            url: "/filterSave",
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
            data: JSON.stringify({ keywords: keyWords, categories: categories })
        })
        .done(function( msg ) {
            if(msg.status == "err"){
                alert(msg.error);
            }
            else{
                window.location.reload();
            }
        })
        .fail(function() {
            alert("Sorry. Server unavailable. ");
        });
    });
    $("#filterResetBtn").click(function(e){
        e.preventDefault();
        $.ajax({
            method: "POST",
            url: "/filterReset",
            contentType: 'application/json; charset=utf-8',
            dataType: 'json'
        })
        .done(function( msg ) {
            if(msg.status == "err"){
                alert(msg.error);
            }
            else{
                window.location.reload();
            }
        })
        .fail(function() {
            alert("Sorry. Server unavailable. ");
            return false;
        });
    });




    function showFormErr(msg){
        $("#formErr").html(msg).show();
    }

    $("#registerBtn").click(function(e){
        e.preventDefault();

        $("#formErr").html("").hide();
        var inputs = $(e.target).parent('form').find(':input');
        inputs.splice(3,1);
        var values = {};
        inputs.each(function() {
            values[this.name] = $(this).val();
        });

        var regex = /^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/;
        if(!regex.test(values.email)){
            showFormErr("Email не валидный");
        }
        else if(values.password == ""){
            showFormErr("Пароль не может быть пустым");
        }
        else if(values.password != values.confirm){
            showFormErr("Пароли не совпадают");
        }
        else{
            $.ajax({
                method: "POST",
                url: "/register/",
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                data: JSON.stringify(values)
            })
            .done(function( msg ) {
                if(msg.status == "err"){
                    showFormErr(msg.error);
                }
                else{
                    window.location = "/confirmMessage/";
                }
            })
            .fail(function() {
                alert("Sorry. Server unavailable. ");
            });
        }
    });
    $("#loginBtn").click(function(e){
        e.preventDefault();

        $("#formErr").html("").hide();
        var inputs = $(e.target).parent('form').find(':input');
        inputs.splice(2,1);
        var values = {};
        inputs.each(function() {
            values[this.name] = $(this).val();
        });

        if(values.email == "" && values.password == ""){
            showFormErr("Поля не могут быть пустыми");
        }
        else{
            $.ajax({
                method: "POST",
                url: "/login/",
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                data: JSON.stringify(values)
            })
            .done(function( msg ) {
                if(msg.status == "err"){
                    showFormErr(msg.error);
                }
                else{
                    window.location = "/dashboard/";
                }
            })
            .fail(function() {
                alert("Sorry. Server unavailable. ");
            });
        }
    });
    $('#notifySound').change(function(e){
        if($(this).is(':checked')){
            notifySound = true;
        }
        else{
            notifySound = false;
        }
    });


    /*********************/
    function siteToIcon(site) {
        var icon = "";
        switch (site) {
            case "f-l":
                icon = "free-lance.ru.gif";
                break;
            case "wl":
                icon = "weblancer.net.gif";
                break;
            case "fl":
                icon = "freelance.png";
                break;
            case "flm":
                icon = "freelancim.png";
                break;
        }
        return icon;
    }
    function siteToName(site) {
        var icon = ""
        switch (site) {
            case "f-l":
                icon = "fl.ru";
                break;
            case "wl":
                icon = "weblancer.net";
                break;
            case "fl":
                icon = "freelance.ru";
                break;
            case "flm":
                icon = "freelansim.ru";
                break;
        }
        return icon;
    }
    function toFullLink(site, href) {
        var link = ""
        switch (site) {
            case "f-l":
                link = "https://www.fl.ru";
                break;
            case "wl":
                link = "https://www.weblancer.net";
                break;
            case "fl":
                link = "freelance.png";
                break;
            case "flm":
                link = "freelancim.png";
                break;
        }
        link += href;
        return link;
    }
    function formatTime(dateStr){
        var d = new Date(dateStr);
        var month = (''+(d.getMonth()+1)).length<2?'0'+(d.getMonth()+1):d.getMonth()+1;
        var day = (''+d.getDate()).length<2?'0'+d.getDate():d.getDate();
        var hours = (''+d.getHours()).length<2?'0'+d.getHours():d.getHours();
        var minutes = (''+d.getMinutes()).length<2?'0'+d.getMinutes():d.getMinutes();
        return day+'.'+month+' '+hours+':'+minutes;
    }
    function addProject(p){
        var projectsTable = $('#projectsList');
        formatTime(p.projectDate);
        $('#projectsList tr:first').before('<tr data-description=\''+p.projectDescription+
            '\'><td><img src="/static/img/'+siteToIcon(p.site)+'" alt=""></td><td><a rel="nofolow" href="'+toFullLink(p.site, p.projectHref)+
            '">'+p.projectTitle+'</a></td><td class="nowrap hidden-xs">'+p.projectPrice+'</td><td class="nowrap hidden-xs col-md-2">'+formatTime(p.projectDate)+'</td></tr>');


        var curTime = Date.now();
        var deltaTime = (curTime - lastNotifySoundTime) / 1000;
        if(deltaTime > 5){
            if(notifySound){
                sound.play();
            }
            lastNotifySoundTime = curTime;
        }

    }
    function getCookie(name) {
        var value = "; " + document.cookie;
        var parts = value.split("; " + name + "=");
        if (parts.length == 2) return parts.pop().split(";").shift();
    }


    if(window.location.pathname == '/dashboard/'){
        var socket = io({path: '/socket.io'});
        socket.on('connect', function() {
            console.log('WID = ', socket.id,  getCookie('sid'))
            var conf = {sid: getCookie('sid')};
            //Session id from cockie
            socket.emit('conn', JSON.stringify(conf));
        });
        socket.on('newProject', function(msg, sendAckCb){
            console.log('msg1', msg)
            addProject(msg)
            sendAckCb("ok");
        });
        socket.on('pingSocket', function(msg, sendAckCb){
            console.log('msg2', msg)
            sendAckCb("ok");
        });
    }

});
