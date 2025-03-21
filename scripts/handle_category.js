let checkboxes = document.querySelectorAll('input[name="category"]');

function getSelectedCategories() {
    var selectedValues = Array.from(checkboxes)
        .filter(checkbox => checkbox.checked)
        .map(checkbox => checkbox.value);

    selectedValues = [...new Set(selectedValues)]
    
    return selectedValues.join(" ")
}
