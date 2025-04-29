// Common functions for board creation and display

// Function to create a grid for a board
function createGrid(containerId, size = 10) {
    const container = $(`#${containerId}`);
    container.empty();
    
    // Create grid header (column numbers)
    const gridHeader = $('<div class="grid-header"></div>');
    gridHeader.append('<div class="grid-header-item"></div>'); // Empty corner cell
    
    for (let i = 1; i <= size; i++) {
        gridHeader.append(`<div class="grid-header-item">${i}</div>`);
    }
    
    container.append(gridHeader);
    
    // Create grid rows
    for (let row = 0; row < size; row++) {
        const gridRow = $('<div class="grid-row"></div>');
        
        // Row header (letters A-J)
        gridRow.append(`<div class="grid-row-header">${String.fromCharCode(65 + row)}</div>`);
        
        // Cells
        for (let col = 0; col < size; col++) {
            gridRow.append(`<div class="cell" data-coord="${String.fromCharCode(65 + row)}${col + 1}"></div>`);
        }
        
        container.append(gridRow);
    }
}

// Store the previous state of each board
const boardStates = {};

// Function to get current board state
function getCurrentBoardState(boardId) {
    const state = {
        cells: {},
        capacity: $(`#${boardId}-capacity`).text()
    };
    
    $(`#${boardId} .cell`).each(function() {
        const $cell = $(this);
        const coord = $cell.data('coord');
        state.cells[coord] = {
            classes: Array.from($cell[0].classList),
            html: $cell.html()
        };
    });
    
    return state;
}

// Function to update a board display
function updateBoard(boardId, playerName, boardData, showPlants = false) {
    // Get current state
    const currentState = boardStates[boardId] || getCurrentBoardState(boardId);
    
    // Create new state
    const newState = {
        cells: {},
        capacity: `Capacity: ${boardData.capacity} / ${boardData.total_capacity}`
    };
    
    // Update capacity only if changed
    if (currentState.capacity !== newState.capacity) {
        $(`#${boardId}-name`).text(playerName);
        $(`#${boardId}-capacity`).text(newState.capacity);
    }
    
    // Process hits and misses first
    if (boardData.hits) {
        boardData.hits.forEach(coord => {
            newState.cells[coord] = {
                classes: ['cell', 'hit'],
                html: ''
            };
        });
    }
    
    if (boardData.misses) {
        boardData.misses.forEach(coord => {
            newState.cells[coord] = {
                classes: ['cell', 'miss'],
                html: ''
            };
        });
    }
    
    // Process plants if they should be shown
    if (showPlants && boardData.plants) {
        boardData.plants.forEach(plant => {
            const plantType = plant.type.toLowerCase();
            const plantCoords = plant.coordinates;
            
            plantCoords.forEach(coord => {
                const isHit = boardData.hits && plantCoords.some(pc => boardData.hits.includes(pc));
                newState.cells[coord] = {
                    classes: ['cell', 'plant-container', 'normal', isHit ? 'plant-damaged' : 'plant-working'],
                    html: `<img src="assets/img/${plantType}.png" alt="${plantType}">`
                };
            });
        });
    }
    
    // Update only cells that have changed
    Object.keys(newState.cells).forEach(coord => {
        const cell = $(`#${boardId} .cell[data-coord="${coord}"]`);
        const newCellState = newState.cells[coord];
        const currentCellState = currentState.cells[coord] || { classes: ['cell'], html: '' };
        
        // Compare classes
        const currentClasses = currentCellState.classes.sort().join(' ');
        const newClasses = newCellState.classes.sort().join(' ');
        
        if (currentClasses !== newClasses || currentCellState.html !== newCellState.html) {
            cell.attr('class', newClasses).html(newCellState.html);
        }
    });
    
    // Clear cells that are not in the new state
    Object.keys(currentState.cells).forEach(coord => {
        if (!newState.cells[coord]) {
            const cell = $(`#${boardId} .cell[data-coord="${coord}"]`);
            if (cell.attr('class') !== 'cell') {
                cell.attr('class', 'cell').empty();
            }
        }
    });
    
    // Store new state
    boardStates[boardId] = newState;
}

// Function to create a board container
function createBoardContainer(boardId, playerName) {
    const boardContainer = $(`
        <div class="board">
            <h3 id="${boardId}-name">${playerName}</h3>
            <p id="${boardId}-capacity">Capacity: 0 / 0</p>
            <div id="${boardId}" class="board-grid"></div>
        </div>
    `);
    
    return boardContainer;
}

// Function to fetch board data based on visibility
async function fetchBoardData(gameId, playerName, isPublic) {
    try {
        // For public games, try to get the full board first
        if (isPublic) {
            const response = await $.ajax({
                url: `/api/games/${gameId}/players/${playerName}/board`,
                type: 'GET'
            });
            return response;
        }
        
        // For private games or if full board fetch fails, get the opponent board
        const response = await $.ajax({
            url: `/api/games/${gameId}/opponent/${playerName}/board`,
            type: 'GET'
        });
        return response;
    } catch (error) {
        console.error('Error fetching board data:', error);
        throw error;
    }
}

// Function to show error message
function showError(message) {
    const errorElement = $('#error-message');
    errorElement.text(message);
    errorElement.show();
    $('#success-message').hide();
}

// Function to show success message
function showSuccess(message, isWaiting = false) {
    const successElement = $('#success-message');
    successElement.text(message);
    
    // Change background color for waiting messages
    if (isWaiting) {
        successElement.css('background-color', '#95a5a6'); // Gray for waiting messages
    } else {
        successElement.css('background-color', '#2ecc71'); // Green for normal success messages
    }
    
    successElement.show();
    $('#error-message').hide();
}
