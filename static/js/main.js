
$( document ).ready(function() {

    $("#projectsList tr").click(function(e){
        if(!e.ctrlKey){
            e.preventDefault();
            var row = $(e.target).closest("tr");
            var lnk = row.first().first().find("a");
            var modal = $("#projectModal");
            modal.find(".modal-title").html(lnk.html());
            modal.find(".modal-body").html(row.data("description"));
            modal.find("#toProjectBtn").attr("href", lnk.attr("href"));
            console.log("BB", modal.find("#toProjectBtn"))
            modal.modal('show');
        }
    });
    $("#settingsModalSaveBtn").click(function(e){
        e.preventDefault();
        var categories = [];
        var catContainerElem = $(".categoriesContainer");
        $.each(catContainerElem, function(k, v){
            var checboxs = $(v).find("input:checkbox");
            $.each(checboxs, function(i, c){
                c = $(c);
                if(c.is(":checked")){
                    categories.push(c.val());
                }

            })
        });
        $.ajax({
            method: "POST",
            url: "/saveFilter",
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
            data: JSON.stringify({ categories: categories })
        })
        .done(function( msg ) {
            window.location.reload();
        })
        .fail(function() {
            alert("Sorry. Server unavailable. ");
        });
    });
    $("#keyWordsAcceptBtn").click(function(e){
        e.preventDefault();
        var keyWords = $("#keyWords").tagsinput('items');
        $.ajax({
            method: "POST",
            url: "/saveKeyWords",
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
            data: JSON.stringify({ keywords: keyWords })
        })
        .done(function( msg ) {
            window.location.reload();
        })
        .fail(function() {
            alert("Sorry. Server unavailable. ");
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
});
