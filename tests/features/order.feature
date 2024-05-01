Feature: order
    In order to manage the orders
    As an user
    I want be able to manage an order

    Scenario: Create an order with an item
        Given I create an order
        And I added an item to the order
        When I retrieve the order
        Then the order should have the item
        And the order state should be "Created"