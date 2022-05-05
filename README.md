# RFC

1. Ustawianie nazw użytkowników
2. Podstawowy JavaScriptowy klient
    - wyświetlanie przychodzących wiadomości
    - potwierdzenia wykonania operacji
    - obsługa błędów
3. Unit testy
4. Testy integracyjne
    - możliwość użycia podstawowego klienta z RFC 2, środowisko klienta na Node.js
    - skrypt bashowy uruchamiający serwer, klienta i testy
5. Refactoring i stworzenie struktury przechowującej stan serwera zamiast zmiennych pakietu `awesomeProjectSlack/internal/handlers`
6. Dodanie trwałości dla przechowywania historii kanałów i wiadomości
    - dodatkowa operacja do pobrania historii wiadomości na kanale z paginacją
7. Logowanie błędów i trwałe ich przechowywanie na zewnętrznym serwerze oddzielnie od historii kanałów i wiadomości
8. Zahostowanie prototypu aplikacji w chmurze AWS
9. Dodanie autentykacji za pomocą oddzielnego serwisu (OAuth2)
    - zarządzanie kontami użytkowników
10. Stworzenie gatewaya będącego warstwą pośredniczącą przed połączeniem się z serwerem czatu
    - throttling na przychodzących połączeniach
    - użycie kolejki (np. AWS SQS bo jest SaaSem i sam się skaluje) dla przychodzących połączeń i konsumowanie ich w mniejszych porcjach na wypadek wysokiego obciążenia serwera
11. Kanały z limitem uczestników
12. Kanały prywatne tylko dla zaproszonych użytkowników
