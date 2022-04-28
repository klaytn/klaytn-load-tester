pragma solidity ^0.4.0;

contract CPUHeavy {
    event finish(uint size, uint signature);
    uint counter;

    function sort(uint size, uint signature) public {
        uint[] memory data = new uint[](size);
        for (uint x = 0; x < data.length; x++) {
            data[x] = size-x;
        }
        quickSort(data, 0, data.length - 1);
        emit finish(size, signature);
    }

    function quickSort(uint[] arr, uint left, uint right) internal {
        uint i;
        uint j;
        uint pivot;

        if (left < right){
            pivot = left;
            i = left;
            j = right;

            while (i < j) {
                while (arr[i] <= arr[pivot] && i < right) i++;
                while (arr[j] > arr[pivot]) j--;
                if (i < j) {
                    (arr[i], arr[j]) = (arr[j], arr[i]);
                }
            }

            (arr[pivot], arr[j]) = (arr[j], arr[pivot]);

            if (j > 1) quickSort(arr, left, j-1);
            quickSort(arr, j+1, right);
        }
    }

    function  empty() public {
        counter++;
    }

    // multinode tester
    uint[20] data20;

    function sortSingle(uint signature) public {
        uint size = 20;
        for (uint x = 0; x < data20.length; x++) {
            data20[x] = size-x;
        }
        quickSortSingle(0, data20.length - 1);
        emit finish(size, signature);
    }

    function checkResult() public view returns(bool) {
        uint prev = data20[0];
        uint cur;
        for (uint x = 1; x < data20.length; x++) {
            cur = data20[x]; 
            if (prev > cur) {
                return false;
            }
            prev = cur;
        }
        return true;
    }

    function quickSortSingle(uint left, uint right) internal {
        uint i;
        uint j;
        uint pivot;

        if (left < right){
            pivot = left;
            i = left;
            j = right;

            while (i < j) {
                while (data20[i] <= data20[pivot] && i < right) i++;
                while (data20[j] > data20[pivot]) j--;
                if (i < j) {
                    (data20[i], data20[j]) = (data20[j], data20[i]);
                }
            }

            (data20[pivot], data20[j]) = (data20[j], data20[pivot]);

            if (j > 1) quickSortSingle(left, j-1);
            quickSortSingle(j+1, right);
        }
    }
}
