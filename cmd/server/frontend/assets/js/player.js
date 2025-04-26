$(document).ready(function() {
    // Get game ID, player name, and token from URL parameters
    const urlParams = new URLSearchParams(window.location.search);
    const gameId = urlParams.get('id');
    const playerName = urlParams.get('player');
    const token = urlParams.get('token');
    
    // Store player name in window object for reference in other functions
    window.playerName = playerName;
    
    if (!gameId || !playerName) {
        showError('Game ID and player name are required!');
        return;
    }
    
    // Update game ID display
    $('#game-id').text(`Game ID: ${gameId}`);
    
    // Game state variables
    let gameData = null;
    let opponentName = null;
    let selectedPlantType = null;
    let boardSize = 10; // Default size
    let requiredCapacity = 0;
    let currentCapacity = 0;
    let playerBoard = {
        plants: [],
        hits: [],
        misses: [],
        total_capacity: 0,
        capacity: 0
    };
    let opponentBoard = {
        hits: [],
        misses: [],
        total_capacity: 0,
        capacity: 0
    };
    let continuePolling = true;
    let isPlayerReady = false;
    
    // Create initial grids
    createGrid('player-board', boardSize);
    createGrid('opponent-board', boardSize);
    
    // Function to create a grid
    function createGrid(containerId, size) {
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
                const coord = `${String.fromCharCode(65 + row)}${col + 1}`;
                gridRow.append(`<div class="cell" data-coord="${coord}"></div>`);
            }
            
            container.append(gridRow);
        }
    }
    
    // Function to update the player's board display
    function updatePlayerBoard(forceRefresh = false) {
        // Only refresh the board if the game is in progress, if forceRefresh is true, or if the player is ready
        if ((gameData && gameData.status !== 'PENDING') || forceRefresh || isPlayerReady) {
            // Clear all cell classes and images
            $('#player-board .cell').removeClass('hit miss normal plant-container plant-working plant-damaged').empty();
            
            // Update capacity display
            $('#player-capacity').text(`Capacity: ${playerBoard.capacity} / ${playerBoard.total_capacity}`);
            
            // Create a map of coordinates to plant types
            const plantMap = {};
            const plantCells = {};
            
            // Ensure plants is initialized
            if (!playerBoard.plants) {
                playerBoard.plants = [];
            }
            
            // Process plants
            playerBoard.plants.forEach(plant => {
                const plantType = plant.type.toLowerCase();
                const plantCoords = [];
                
                plant.coordinates.forEach(coord => {
                    plantMap[coord] = plantType;
                    plantCoords.push(coord);
                });
                
                plantCells[plantType] = plantCells[plantType] || [];
                plantCells[plantType].push(plantCoords);
            });
            
            // First, add all plant images and containers
            for (const [coord, plantType] of Object.entries(plantMap)) {
                const cell = $(`#player-board .cell[data-coord="${coord}"]`);
                if (cell.length) {
                    cell.html(`<img src="assets/img/${plantType}.png" alt="${plantType}">`);
                    cell.addClass('plant-container normal');
                }
            }
            
            // Then mark hits and misses
            if (playerBoard.hits) {
                playerBoard.hits.forEach(coord => {
                    const cell = $(`#player-board .cell[data-coord="${coord}"]`);
                    if (cell.length) {
                        cell.removeClass('normal').addClass('hit');
                    }
                });
            }
            
            if (playerBoard.misses) {
                playerBoard.misses.forEach(coord => {
                    const cell = $(`#player-board .cell[data-coord="${coord}"]`);
                    if (cell.length) {
                        cell.removeClass('normal').addClass('miss');
                    }
                });
            }
            
            // Finally, add working/damaged status to plants
            for (const [plantType, plantGroups] of Object.entries(plantCells)) {
                plantGroups.forEach(plantCoords => {
                    // Check if any cell in this plant is hit
                    const isHit = playerBoard.hits && plantCoords.some(coord => playerBoard.hits.includes(coord));
                    
                    // Add status classes to all cells in this plant
                    plantCoords.forEach(coord => {
                        const cell = $(`#player-board .cell[data-coord="${coord}"]`);
                        if (cell.length) {
                            if (isHit) {
                                cell.addClass('plant-damaged');
                            } else {
                                cell.addClass('plant-working');
                            }
                        }
                    });
                });
            }
        }
    }
    
    // Function to update the opponent's board display
    function updateOpponentBoard() {
        // Clear all cell classes and images
        $('#opponent-board .cell').removeClass('hit miss normal').empty();
        
        // Update capacity display
        $('#opponent-capacity').text(`Capacity: ${opponentBoard.capacity} / ${opponentBoard.total_capacity}`);
        
        // Mark hits and misses
        if (opponentBoard.hits) {
            opponentBoard.hits.forEach(coord => {
                const cell = $(`#opponent-board .cell[data-coord="${coord}"]`);
                if (cell.length) {
                    cell.addClass('hit');
                }
            });
        }
        
        if (opponentBoard.misses) {
            opponentBoard.misses.forEach(coord => {
                const cell = $(`#opponent-board .cell[data-coord="${coord}"]`);
                if (cell.length) {
                    cell.addClass('miss');
                }
            });
        }
    }
    
    // Function to create a board container for an opponent
    function createOpponentBoardContainer(opponentName, index) {
        const boardId = `opponent${index}-board`;
        const boardContainer = $(`
            <div class="board">
                <h3>${opponentName}'s Board</h3>
                <p id="opponent${index}-capacity">Capacity: 0 / 0</p>
                <div id="${boardId}"></div>
            </div>
        `);
        
        return { boardContainer, boardId };
    }
    
    // Function to fetch game data and update the UI
    function updateGameData() {
        if (!continuePolling) return;
        
        $.ajax({
            url: `/api/games/${gameId}`,
            type: 'GET',
            success: function(data) {
                gameData = data;
                
                // Update game status information
                $('#game-status').text(`Status: ${data.status}`);
                $('#game-turn').text(`Turn: ${data.turn || 'N/A'}`);
                $('#game-winner').text(`Winner: ${data.winner || 'None'}`);
                
                // Get player names
                const playerNames = Object.keys(data.players);
                
                // Update players list
                const playersList = $('#players');
                playersList.empty();
                
                playerNames.forEach(playerName => {
                    const playerInfo = data.players[playerName];
                    const readyStatus = playerInfo.ready ? 'Ready' : 'Not Ready';
                    const isCurrentPlayer = playerName === window.playerName;
                    const playerLabel = isCurrentPlayer ? `${playerName} (You)` : playerName;
                    playersList.append(`<li>${playerLabel} - ${readyStatus} (Capacity: ${playerInfo.capacity}/${playerInfo.total_capacity})</li>`);
                });
                
                // Find opponents if there are at least 2 players
                if (playerNames.length >= 2) {
                    const opponents = playerNames.filter(name => name !== playerName);
                    
                    // Update player info
                    const playerInfo = data.players[playerName];
                    
                    // Update capacity requirements
                    requiredCapacity = data.capacity || 1000; // Default to 1000 if not provided
                    $('#required-capacity').text(`Required: ${requiredCapacity}`);
                    
                    // Check if player is ready
                    isPlayerReady = playerInfo.ready;
                    
                    // Update UI based on game status
                    updateUIForGameStatus(data.status, data.turn);
                    
                    // Fetch player board
                    fetchPlayerBoard();
                    
                    // Clear opponents container
                    const opponentsContainer = $('#opponents-container');
                    opponentsContainer.empty();
                    
                    // Create a board for each opponent
                    opponents.forEach((opponent, index) => {
                        const { boardContainer, boardId } = createOpponentBoardContainer(opponent, index + 1);
                        opponentsContainer.append(boardContainer);
                        
                        // Create the grid for this board
                        createGrid(boardId);
                        
                        // Fetch board data for this opponent
                        fetchOpponentBoard(opponent, boardId);
                    });
                } else {
                    // Only one player (current player) in the game
                    $('#opponents-container').html(`
                        <div class="board">
                            <h3>Waiting for Opponents</h3>
                            <p>No opponents have joined yet.</p>
                        </div>
                    `);
                }
            },
            error: function(xhr) {
                console.error('Error fetching game data:', xhr);
                
                // Check if the game doesn't exist
                if (xhr.status === 404) {
                    showError('Game not found. Please check the game ID and try again.');
                    continuePolling = false; // Stop polling
                    
                    // Hide setup controls and boards
                    $('#setup-controls').hide();
                    $('.boards-container').hide();
                    
                    // Show a button to go back to home
                    $('.home-link').show();
                } else {
                    showError('Error fetching game data. Please try again.');
                }
            }
        });
    }
    
    // Function to update UI based on game status
    function updateUIForGameStatus(status, turn) {
        if (status === 'PENDING') {
            // Show setup controls if player is not ready
            if (!isPlayerReady) {
                // Show setup controls and hide strike controls
                $('#setup-controls').show();
                $('#strike-controls').hide();
                
                // Move setup controls next to the player's board
                $('#setup-controls').insertAfter('.board:first-child');
                
                // Hide opponents container completely
                $('#opponents-container').hide();
            } else {
                // Player is ready but game hasn't started yet
                $('#setup-controls').hide();
                $('#strike-controls').hide();
                
                // Show waiting message in opponents container
                $('#opponents-container').show().html(`
                    <div class="board">
                        <h3>Waiting for Opponents</h3>
                        <p>Waiting for other players to be ready...</p>
                    </div>
                `);
                
                showSuccess('You are ready! Waiting for other players to be ready...');
            }
        } else if (status === 'IN_PROGRESS') {
            // Hide setup controls
            $('#setup-controls').hide();
            
            // Hide all setup-related elements
            $('.plant-selector').hide();
            $('#orientation-selector').hide();
            $('.instructions').hide();
            $('.actions').hide();
            
            // Always show strike controls during gameplay
            $('#strike-controls').show();
            
            // Make sure opponents container is visible
            $('#opponents-container').show();
            
            // Make sure opponent boards are visible
            const opponents = Object.keys(gameData.players).filter(name => name !== playerName);
            if (opponents.length > 0) {
                // Fetch boards for all opponents
                opponents.forEach((opponent, index) => {
                    fetchOpponentBoard(opponent, `opponent${index+1}-board`);
                });
            }
            
            // Update status message based on turn
            if (turn === playerName) {
                showSuccess('It\'s your turn! Click on the opponent\'s board to strike.');
            } else {
                showSuccess('Waiting for opponent\'s move...', true); // true indicates it's a waiting message
            }
        } else if (status === 'END') {
            // Game is over
            $('#setup-controls').hide();
            $('#strike-controls').hide();
            
            // Hide all setup-related elements
            $('.plant-selector').hide();
            $('#orientation-selector').hide();
            $('.instructions').hide();
            $('.actions').hide();
            
            if (gameData.winner === playerName) {
                showSuccess('Congratulations! You won the game!');
            } else {
                showError('Game over. You lost the game.');
            }
        }
    }
    
    // Function to fetch player's board
    function fetchPlayerBoard() {
        // Skip fetching if we have unsaved local changes
        if (playerBoard.plants && playerBoard.plants.length > 0 && !isPlayerReady) {
            return;
        }
        
        $.ajax({
            url: `/api/games/${gameId}/players/${playerName}/board`,
            type: 'GET',
            success: function(data) {
                playerBoard = data;
                // Ensure plants is initialized
                if (!playerBoard.plants) {
                    playerBoard.plants = [];
                }
                updatePlayerBoard();
                
                // Update current capacity
                currentCapacity = data.total_capacity;
                $('#current-capacity').text(`Current Capacity: ${currentCapacity}`);
                
                // Enable/disable ready button based on capacity
                updateReadyButton();
            },
            error: function(xhr) {
                // If board not found, it might not be set yet
                console.log('Player board not set yet');
                // Make sure playerBoard.plants is initialized
                if (!playerBoard.plants) {
                    playerBoard.plants = [];
                }
            }
        });
    }
    
    // Function to fetch opponent's board
    function fetchOpponentBoard(opponentName, boardId) {
        if (!opponentName) return;
        
        $.ajax({
            url: `/api/games/${gameId}/opponent/${opponentName}/board`,
            type: 'GET',
            success: function(data) {
                // Update the board display
                updateOpponentBoardDisplay(boardId, opponentName, data);
                
                // Store the data for the current opponent if it's the first one
                if (boardId === 'opponent1-board') {
                    opponentBoard = data;
                }
            },
            error: function(xhr) {
                console.error(`Error fetching board for ${opponentName}:`, xhr);
                // Even if there's an error, we'll try again on the next polling cycle
            }
        });
    }
    
    // Function to update an opponent's board display
    function updateOpponentBoardDisplay(boardId, opponentName, board) {
        // Extract player ID from board ID
        const playerId = boardId.split('-')[0];
        
        // Update capacity display
        $(`#${playerId}-capacity`).text(`Capacity: ${board.capacity} / ${board.total_capacity}`);
        
        // Check if the grid exists, if not create it
        if ($(`#${boardId} .grid-header`).length === 0) {
            console.log(`Creating grid for ${boardId}`);
            createGrid(boardId);
        }
        
        // Clear all cell classes
        $(`#${boardId} .cell`).removeClass('hit miss').empty();
        
        // Mark hits and misses
        if (board.hits) {
            board.hits.forEach(coord => {
                const cell = $(`#${boardId} .cell[data-coord="${coord}"]`);
                if (cell.length) {
                    cell.addClass('hit');
                }
            });
        }
        
        if (board.misses) {
            board.misses.forEach(coord => {
                const cell = $(`#${boardId} .cell[data-coord="${coord}"]`);
                if (cell.length) {
                    cell.addClass('miss');
                }
            });
        }
    }
    
    // Function to update the ready button state
    function updateReadyButton() {
        const readyBtn = $('#ready-btn');
        
        // Check if the current capacity meets the requirements
        // The capacity should be at least the required capacity and at most 10% extra
        console.log(`Current capacity: ${currentCapacity}, Required: ${requiredCapacity}, Max: ${requiredCapacity * 1.1}`);
        
        // Debug: Force enable the button if capacity is at least 1000
        if (currentCapacity >= 1000) {
            console.log("Capacity is at least 1000, enabling Ready button");
            readyBtn.removeClass('button-disabled');
            return;
        }
        
        if (currentCapacity >= requiredCapacity && currentCapacity <= requiredCapacity * 1.1) {
            console.log("Capacity meets requirements, enabling Ready button");
            readyBtn.removeClass('button-disabled');
        } else {
            console.log("Capacity does not meet requirements, disabling Ready button");
            readyBtn.addClass('button-disabled');
        }
    }
    
    // Function to place a plant on the board
    function placePlant(coord, type) {
        if (!type) return;
        
        // Ensure playerBoard.plants is initialized
        if (!playerBoard.plants) {
            playerBoard.plants = [];
        }
        
        // Get plant size
        let width, height;
        switch (type) {
            case 'NUCLEAR':
                width = height = 3;
                break;
            case 'GAS':
                width = height = 2;
                break;
            case 'WIND':
                if ($('input[name="orientation"]:checked').val() === 'horizontal') {
                    width = 2;
                    height = 1;
                } else {
                    width = 1;
                    height = 2;
                }
                break;
            case 'SOLAR':
                width = height = 1;
                break;
            default:
                return;
        }
        
        // Parse coordinate
        const row = coord.charCodeAt(0) - 65; // A=0, B=1, etc.
        const col = parseInt(coord.substring(1)) - 1; // 1-indexed to 0-indexed
        
        // Check if plant fits on the board
        if (row + height > boardSize || col + width > boardSize) {
            showError('Plant does not fit on the board!');
            return false;
        }
        
        // Generate all coordinates for the plant
        const coordinates = [];
        for (let r = 0; r < height; r++) {
            for (let c = 0; c < width; c++) {
                const newRow = row + r;
                const newCol = col + c;
                const newCoord = `${String.fromCharCode(65 + newRow)}${newCol + 1}`;
                coordinates.push(newCoord);
            }
        }
        
        // Check if any of the coordinates are already occupied
        if (playerBoard.plants && playerBoard.plants.length > 0) {
            for (const plantCoord of coordinates) {
                for (const plant of playerBoard.plants) {
                    if (plant.coordinates.includes(plantCoord)) {
                        showError('Cannot place plant here. Space already occupied!');
                        return false;
                    }
                }
            }
        }
        
        // Add the plant to the board
        playerBoard.plants.push({
            type: type,
            coordinates: coordinates
        });
        
        // Update capacity
        let plantCapacity;
        switch (type) {
            case 'NUCLEAR':
                plantCapacity = 1000;
                break;
            case 'GAS':
                plantCapacity = 300;
                break;
            case 'WIND':
                plantCapacity = 100;
                break;
            case 'SOLAR':
                plantCapacity = 25;
                break;
            default:
                plantCapacity = 0;
        }
        
        // Update capacity values
        currentCapacity += plantCapacity;
        playerBoard.total_capacity = currentCapacity;
        playerBoard.capacity = currentCapacity;
        
        console.log(`Plant placed: ${type}, Capacity added: ${plantCapacity}, New total: ${currentCapacity}`);
        
        // Directly update the DOM with the plant image
        for (const coord of coordinates) {
            const cell = $(`#player-board .cell[data-coord="${coord}"]`);
            if (cell.length) {
                cell.html(`<img src="assets/img/${type.toLowerCase()}.png" alt="${type}">`);
                cell.addClass('plant-container normal plant-working');
            }
        }
        
        // Update capacity display
        $('#current-capacity').text(`Current Capacity: ${currentCapacity}`);
        
        // Update ready button state
        updateReadyButton();
        
        return true;
    }
    
    // Function to save the board to the server
    function saveBoard(callback) {
        // Ensure playerBoard.plants is initialized
        if (!playerBoard.plants) {
            playerBoard.plants = [];
        }
        
        $.ajax({
            url: `/api/games/${gameId}/players/${playerName}/board`,
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                plants: playerBoard.plants
            }),
            success: function(data) {
                console.log('Board saved successfully:', data);
                playerBoard = data;
                updatePlayerBoard();
                
                // Call the callback function if provided
                if (callback && typeof callback === 'function') {
                    callback();
                }
            },
            error: function(xhr) {
                console.error('Error saving board:', xhr);
                showError('Error saving board: ' + (xhr.responseJSON ? xhr.responseJSON.error : 'Unknown error'));
            }
        });
    }
    
    // Function to mark player as ready
    function setPlayerReady() {
        $.ajax({
            url: `/api/games/${gameId}/players/${playerName}/ready`,
            type: 'POST',
            success: function(data) {
                console.log('Player marked as ready:', data);
                isPlayerReady = true;
                showSuccess('You are ready! Waiting for opponent...');
                $('#setup-controls').hide();
                $('#opponents-container').show();
                console.log("Player is ready, hiding setup controls");
            },
            error: function(xhr) {
                console.error('Error marking player as ready:', xhr);
                showError('Error marking player as ready: ' + (xhr.responseJSON ? xhr.responseJSON.error : 'Unknown error'));
            }
        });
    }
    
    // Function to strike opponent's board
    function strikeOpponent(coord, targetOpponent) {
        if (!gameData || gameData.status !== 'IN_PROGRESS' || gameData.turn !== playerName) {
            showError('It\'s not your turn to strike!');
            return;
        }
        
        // Use the provided target opponent or fall back to the first opponent
        const target = targetOpponent || opponentName;
        
        if (!target) {
            showError('No opponent to strike!');
            return;
        }
        
        // Parse coordinate to get y and x
        const y = coord.charAt(0);
        const x = coord.substring(1);
        
        $.ajax({
            url: `/api/games/${gameId}/players/${playerName}/strike?target=${target}&y=${y}&x=${x}`,
            type: 'POST',
            success: function(data) {
                console.log('Strike result:', data);
                
                if (data.result === 'HIT') {
                    showSuccess(`HIT! You struck ${target}'s power plant!`);
                } else {
                    showSuccess(`MISS! You didn't hit anything on ${target}'s board.`);
                }
                
                // Update game data
                updateGameData();
            },
            error: function(xhr) {
                console.error('Error striking opponent:', xhr);
                showError('Error striking opponent: ' + (xhr.responseJSON ? xhr.responseJSON.error : 'Unknown error'));
            }
        });
    }
    
    // Function to show error message
    function showError(message) {
        const errorElement = $('#error-message');
        errorElement.text(message);
        errorElement.show();
        $('#success-message').hide();
        
        // Don't auto-hide messages
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
        
        // Don't auto-hide messages
    }
    
    // Plant selection event
    $('.plant-option').click(function() {
        $('.plant-option').removeClass('selected');
        $(this).addClass('selected');
        selectedPlantType = $(this).data('type');
        
        // Show/hide orientation selector for Wind plants
        if (selectedPlantType === 'WIND') {
            $('#orientation-selector').show();
        } else {
            $('#orientation-selector').hide();
        }
    });
    
    // Player board cell click event (for placing plants)
    $('#player-board').on('click', '.cell', function() {
        if (gameData && gameData.status !== 'PENDING') return;
        if (isPlayerReady) return;
        
        const coord = $(this).data('coord');
        
        if (selectedPlantType) {
            placePlant(coord, selectedPlantType);
        } else {
            showError('Please select a power plant type first!');
        }
    });
    
    // Opponent board cell click event (for striking)
    $('#opponents-container').on('click', '.cell', function() {
        if (!gameData || gameData.status !== 'IN_PROGRESS' || gameData.turn !== playerName) {
            showError('It\'s not your turn to strike!');
            return;
        }
        
        const coord = $(this).data('coord');
        const boardId = $(this).closest('[id^="opponent"]').attr('id');
        const opponentIndex = parseInt(boardId.replace('opponent', '').replace('-board', '')) - 1;
        
        // Get the opponents list
        const opponents = Object.keys(gameData.players).filter(name => name !== playerName);
        
        // Get the target opponent name
        const targetOpponent = opponents[opponentIndex];
        
        if (!targetOpponent) {
            showError('Invalid opponent selection');
            return;
        }
        
        // Strike the opponent
        strikeOpponent(coord, targetOpponent);
    });
    
    // Reset board button click event
    $('#reset-board-btn').click(function() {
        playerBoard.plants = [];
        playerBoard.hits = [];
        playerBoard.misses = [];
        playerBoard.total_capacity = 0;
        playerBoard.capacity = 0;
        currentCapacity = 0;
        
        // Clear all cell classes and images
        $('#player-board .cell').removeClass('hit miss normal plant-container plant-working plant-damaged').empty();
        
        $('#current-capacity').text(`Current Capacity: ${currentCapacity}`);
        updateReadyButton();
    });
    
    // Ready button click event
    $('#ready-btn').click(function() {
        console.log("Ready button clicked, currentCapacity:", currentCapacity);
        
        // Force enable if capacity is at least 1000
        if (currentCapacity >= 1000) {
            console.log("Capacity is sufficient, proceeding with ready");
            
            // First save the board, then mark as ready
            saveBoard(function() {
                setPlayerReady();
            });
            return;
        }
        
        if ($(this).hasClass('button-disabled')) {
            showError('Your board does not meet the capacity requirements!');
            return;
        }
        
        // First save the board, then mark as ready
        // Use the callback to ensure setPlayerReady is called only after saveBoard completes
        saveBoard(function() {
            setPlayerReady();
        });
    });
    
    // Hide orientation selector initially (only show for Wind plants)
    $('#orientation-selector').hide();
    
    // Initial update
    updateGameData();
    
    // Set up polling
    let pollingInterval;
    
    function startPolling() {
        // Clear any existing interval
        if (pollingInterval) {
            clearInterval(pollingInterval);
        }
        
        // Determine polling frequency based on game status
        const pollingFrequency = (gameData && gameData.status === 'PENDING') ? 3000 : 1000; // 3 seconds for PENDING, 1 second for IN_PROGRESS
        
        pollingInterval = setInterval(function() {
            if (!continuePolling) {
                clearInterval(pollingInterval);
                return;
            }
            updateGameData();
            
            // Adjust polling frequency if game status changes
            if (gameData) {
                const newFrequency = (gameData.status === 'PENDING') ? 3000 : 1000;
                if (newFrequency !== pollingFrequency) {
                    startPolling(); // Restart polling with new frequency
                }
            }
        }, pollingFrequency);
    }
    
    // Initial update and start polling
    updateGameData();
    startPolling();
});
