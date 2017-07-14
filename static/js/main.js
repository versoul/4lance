
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
});
